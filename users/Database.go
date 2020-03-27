package main

import (
	"crypto/tls"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"golang.org/x/net/context"
)

var Wrapper ArangoWrapper

func (wrapper *ArangoWrapper) initDriver(uri string, username string, password string) {
	var err error
	wrapper.Connection, err = http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{uri},
		TLSConfig: &tls.Config{ /*...*/ },
	})
	if err != nil {
		failOnError(err, "An occured has occured while initiating the connection")
	}
	wrapper.Client, err = driver.NewClient(driver.ClientConfig{
		Connection:     wrapper.Connection,
		Authentication: driver.BasicAuthentication(username, password),
	})
	if err != nil {
		failOnError(err, "An occured has occured while setting the client")
	}

	wrapper.ExecContext = context.Background()
	wrapper.Database, err = wrapper.Client.Database(wrapper.ExecContext, "users")
	if err != nil {
		failOnError(err, "An occured has occured while setting the database")
	}

	// define the edgeCollection to store the edges
	var edgeDefinition driver.EdgeDefinition
	edgeDefinition.Collection = "UserEdges"
	edgeDefinition.From = []string{string(UsersCollection)}
	edgeDefinition.To = []string{string(UsersCollection)}

	var options driver.CreateGraphOptions
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}

	wrapper.Graph, err = wrapper.Database.CreateGraph(wrapper.ExecContext, string(UsersGraph), &options)
	if err != nil {
		failOnError(err, "Couldn't create the graph")
		wrapper.Graph, err = wrapper.Database.Graph(wrapper.ExecContext, string(UsersGraph))
	}

	wrapper.Collection, err = wrapper.Graph.CreateVertexCollection(wrapper.ExecContext, string(UsersCollection))
	if err != nil {
		failOnError(err, "Couldn't create the collection")
		wrapper.Collection, err = wrapper.Graph.VertexCollection(wrapper.ExecContext, string(UsersCollection))
	}
}

func (wrapper *ArangoWrapper) Close() {

}

func (wrapper *ArangoWrapper) Set(data map[string]interface{}) bool {
	return false
}

func (wrapper *ArangoWrapper) Create(data interface{}) {

	meta, err := wrapper.Collection.CreateDocument(wrapper.ExecContext, data)
	if err != nil {
		failOnError(err, "An error has occured while trying to add data")
	}
	fmt.Printf("Created document with key '%s', revision '%s'\n", meta.Key, meta.Rev)
}

func (wrapper *ArangoWrapper) AddRelation(collection ArangoCollections, edges ArangoEdge, data RelationModel) bool {
	data.From = fmt.Sprintf("%s/%s", string(collection), data.From)
	data.To = fmt.Sprintf("%s/%s", string(collection), data.To)
	data.Key = fmt.Sprintf("%s_%s_%s", data.From, data.To, string(data.Relation))
	edgeCollection, _, err := wrapper.Graph.EdgeCollection(wrapper.ExecContext, string(edges))
	if err != nil {
		failOnError(err, "An error has occured while trying to select edge")
	}

	_, err = edgeCollection.CreateDocument(wrapper.ExecContext, data)

	if err != nil {
		failOnError(err, "An error has occured while trying to add data")
		return false
	}

	return true
}

func (wrapper *ArangoWrapper) UpdateRelation(collection ArangoCollections, edges ArangoEdge, data RelationModel) bool {
	data.From = fmt.Sprintf("%s/%s", string(collection), data.From)
	data.To = fmt.Sprintf("%s/%s", string(collection), data.To)

	edgeCollection, _, err := wrapper.Graph.EdgeCollection(wrapper.ExecContext, string(edges))
	if err != nil {
		failOnError(err, "An error has occured while trying to select edge")
	}

	_, err = edgeCollection.UpdateDocument(wrapper.ExecContext, fmt.Sprintf("%s_%s_%s", data.From, data.To, string(data.Relation)), data)

	if err != nil {
		failOnError(err, "An error has occured while trying to update data")
		return false
	}

	return true
}

func (wrapper *ArangoWrapper) RemoveRelation(collection ArangoCollections, edges ArangoEdge, data RelationModel) bool {
	data.From = fmt.Sprintf("%s/%s", string(collection), data.From)
	data.To = fmt.Sprintf("%s/%s", string(collection), data.To)

	edgeCollection, _, err := wrapper.Graph.EdgeCollection(wrapper.ExecContext, string(edges))
	if err != nil {
		failOnError(err, "An error has occured while trying to select edge")
	}

	_, err = edgeCollection.RemoveDocument(wrapper.ExecContext, fmt.Sprintf("%s_%s_%s", data.From, data.To, string(data.Relation)))

	if err != nil {
		failOnError(err, "An error has occured while trying to add data")
		return false
	}

	return true
}

func (wrapper *ArangoWrapper) GetWhere(searchKey string, operation ArangOperator, collection ArangoCollections, value string) ([]interface{}, int) {
	data := constructArangoRequest(collection, value)
	query := fmt.Sprintf("FOR data IN @@collection FILTER data.%s %s @value RETURN data", searchKey, string(operation))

	err := wrapper.Database.ValidateQuery(nil, query)
	if err != nil {
		failOnError(err, "An error has occured while validating the query")
	}

	cursor, err := wrapper.Database.Query(wrapper.ExecContext, query, data) //, data
	if err != nil {
		failOnError(err, "An error has occured while querying")
	}
	defer cursor.Close()

	var results []interface{}
	for {
		var result interface{}
		_, err := cursor.ReadDocument(wrapper.ExecContext, &result)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			// handle other errors
		}
		results = append(results, result)
	}

	return results, len(results)
}

func (wrapper *ArangoWrapper) GetConnected(graph ArangoGraph, collection ArangoCollections, direction EdgeDirection, value string, depth int) ([]interface{}, int) {
	data := constructALQRequest(graph, fmt.Sprintf("'%s/%s'", collection, value))
	query := fmt.Sprintf("FOR v, e IN 1..%d %s @value GRAPH @@graph RETURN {user: v, connection: e}", depth, string(direction))

	err := wrapper.Database.ValidateQuery(nil, query)
	if err != nil {
		failOnError(err, "An error has occured while validating the query")
	}

	cursor, err := wrapper.Database.Query(wrapper.ExecContext, query, data) //, data
	if err != nil {
		failOnError(err, "An error has occured while querying")
	}
	defer cursor.Close()

	var results []interface{}
	for {
		var result interface{}
		_, err := cursor.ReadDocument(wrapper.ExecContext, &result)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			// handle other errors
		}
		results = append(results, result)
	}

	return results, len(results)
}

func constructArangoRequest(collection ArangoCollections, value string) map[string]interface{} {
	var data map[string]interface{} = make(map[string]interface{})
	data["value"] = value
	data["@collection"] = string(collection)
	return data
}

func constructALQRequest(graph ArangoGraph, value string) map[string]interface{} {
	var data map[string]interface{} = make(map[string]interface{})
	data["value"] = value
	data["@graph"] = string(graph)
	return data
}

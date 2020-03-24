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
	edgeDefinition.From = []string{"UserCollection"}
	edgeDefinition.To = []string{"UserCollection"}

	var options driver.CreateGraphOptions
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}

	wrapper.Graph, err = wrapper.Database.CreateGraph(wrapper.ExecContext, "users-graph", &options)
	if err != nil {
		failOnError(err, "Couldn't create the graph")
		wrapper.Graph, err = wrapper.Database.Graph(wrapper.ExecContext, "users-graph")
	}

	wrapper.Collection, err = wrapper.Graph.CreateVertexCollection(wrapper.ExecContext, "UserCollection")
	if err != nil {
		failOnError(err, "Couldn't create the collection")
		wrapper.Collection, err = wrapper.Graph.VertexCollection(wrapper.ExecContext, "UserCollection")
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

func (wrapper *ArangoWrapper) SetRelation(data RelationModel) bool {
	data.From = fmt.Sprintf("UserCollection/%s", data.From)
	data.To = fmt.Sprintf("UserCollection/%s", data.To)

	edgeCollection, _, err := wrapper.Graph.EdgeCollection(wrapper.ExecContext, "UserEdges")
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

func (wrapper *ArangoWrapper) GetWhere(value string) ([]interface{}, int) {
	data := constructArangoRequest(value)
	query := "FOR data IN UserCollection FILTER data._key == @value RETURN data"

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

func (wrapper *ArangoWrapper) GetConnected(value string) ([]interface{}, int) {
	data := constructArangoRequest(fmt.Sprintf("'UserCollection/%s'", value))
	query := "FOR v, e IN 1..1 OUTBOUND @value GRAPH 'users-graph' RETURN {user: v, connection: e}"

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

/*
db._query(`FOR v, e IN 1..3 OUTBOUND 'persons/eve'
           GRAPH 'knows_graph'
           RETURN {v: v, e: e}`)
*/

func constructArangoRequest(value string) map[string]interface{} {
	var data map[string]interface{} = make(map[string]interface{})
	data["value"] = value
	return data
}

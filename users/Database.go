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
	edgeDefinition.Collection = "users-edges"
	edgeDefinition.From = []string{"users-collections"}
	edgeDefinition.To = []string{"users-collections"}

	var options driver.CreateGraphOptions
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}

	wrapper.Graph, err = wrapper.Database.CreateGraph(wrapper.ExecContext, "users-graph", &options)
	if err != nil {
		failOnError(err, "Couldn't create the graph")
		wrapper.Graph, err = wrapper.Database.Graph(wrapper.ExecContext, "users-graph")
	}

	wrapper.Collection, err = wrapper.Graph.CreateVertexCollection(wrapper.ExecContext, "users-collections")
	if err != nil {
		failOnError(err, "Couldn't create the collection")
		wrapper.Collection, err = wrapper.Graph.VertexCollection(wrapper.ExecContext, "users-collections")
	}
}

func (wrapper *ArangoWrapper) Close() {

}

func (wrapper *ArangoWrapper) Set(data map[string]interface{}) bool {
	return false
}

func (wrapper *ArangoWrapper) Create(data interface{}) bool {

	meta, err := wrapper.Collection.CreateDocument(wrapper.ExecContext, data)
	if err != nil {
		failOnError(err, "An error has occured while trying to add data")
		return false
	}
	fmt.Printf("Created document with key '%s', revision '%s'\n", meta.Key, meta.Rev)
	return true
}

func (wrapper *ArangoWrapper) SetRelation(data RelationModel) bool {
	data.From = fmt.Sprintf("users-collections/%s", data.From)
	data.To = fmt.Sprintf("users-collections/%s", data.To)

	edgeCollection, _, err := wrapper.Graph.EdgeCollection(wrapper.ExecContext, "users-edges")
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

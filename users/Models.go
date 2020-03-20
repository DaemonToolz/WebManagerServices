package main

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"

	"github.com/arangodb/go-driver"
)

/*
	-------------- Constants
*/

func ObjectToMap(data interface{}) map[string]interface{} {
	b, _ := json.Marshal(data)
	var output map[string]interface{}
	json.Unmarshal(b, &output)
	return output
}

type Relation string

const (
	RELATION_FRIEND Relation = "FRIEND"
	RELATION_SHARE  Relation = "SHARE"
	RELATION_FOLLOW Relation = "FOLLOW"
)

type RelationModel struct {
	From     string `json:"_from"`
	To       string `json:"_to"`
	Relation string `json:"relation"`
}

type ArangoWrapper struct {
	ExecContext context.Context
	Client      driver.Client
	Connection  driver.Connection
	Database    driver.Database
	Graph       driver.Graph
	Collection  driver.Collection
}

type UserInfo struct {
	Username  string    `json:"_key"`
	CreatedAt time.Time `json:"created"`
	RealName  string    `json:"real_name"`
	Email     string    `json:"email"`
}

//	CreatedAt time.Time `json:"created_at"`

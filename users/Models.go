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

type RabbitMqMsg struct {
	ID       string    `json:"id"`
	Date     time.Time `json:"date"`
	Status   int       `json:"status"` // New = 0, Ongoing = 1, Done = 2, Error = 3...
	Function Function  `json:"function"`
	To       string    `json:"to"`
	Priority int       `json:"priority"` // Critical = 0,
	Type     int       `json:"type"`     // Error, warn
	Payload  string    `json:"payload"`
}

type Function string

const (
	UsersRegistered Function = "system.users.registered"
	UsersValidate   Function = "system.users.validated"
)

const (
	UsersChannel string = "users-notification"
)

const ( // iota is reset to 0
	STATUS_ERROR   = iota // 0
	STATUS_NEW     = iota // 1
	STATUS_ONGOING = iota // 2
	STATUS_DONE    = iota // 3
)

const ( // iota is reset to 0
	PRIORITY_LOW      = iota // 0
	PRIORITY_STD      = iota // 1
	PRIORITY_MEDIUM   = iota // 2
	PRIORITY_HIGH     = iota // 3
	PRIORITY_CRITICAL = iota // 3
)

const ( // iota is reset to 0
	TYPE_INFO    = iota // 0
	TYPE_SUCCESS = iota // 1
	TYPE_WARN    = iota // 2
	TYPE_ERROR   = iota // 3
)

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

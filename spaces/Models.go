package main

import "time"

/*
	--------------- MODELS
*/
type Function string

const (
	MySpaceUpdate   Function = "myspace.space_update"
	MySpaceValidate Function = "myspace.space_validation"
	MySpaceNotify   Function = "myspace.notify_all"
)

const (
	MySpaceGeneralChannel string = "myspace-notification"
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

type FileModel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Type int    `json:"type"`
	Size int64  `json:"size"`
}

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

type UserInitialization struct {
	UserId     string    `json:"userid"`
	InitStatus int       `json:"status"`
	Created    bool      `json:"created"`
	CreatedAt  time.Time `json:"created_at"`
}

package main

import "time"

/*
	--------------- MODELS
*/
type Function string

const (
	FilewatchNotify Function = "filewatch.notify"
	FilewatchSysUpd Function = "filewatch.system_updates"
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

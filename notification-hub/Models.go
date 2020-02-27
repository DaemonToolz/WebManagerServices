package main

import "time"

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
	MySpaceUpdate   Function = "myspace.space_update"
	MySpaceValidate Function = "myspace.space_validation"
	FilewatchNotify Function = "filewatch.notify"
	FilewatchSysUpd Function = "filewatch.system_updates"
)

const (
	MySpaceGeneralChannel string = "myspace-notification"
	BroadcastChannel      string = "general-notification"
)

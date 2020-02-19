package main

type RabbitMqMsg struct {
	ID       string `json:"id"`
	Status   int    `json:"status"` // New = 0, Ongoing = 1, Done = 2, Error = 3...
	Function string `json:"function"`
	To       string `json:"to"`
	Priority int    `json:"priority"` // Critical = 0,
	Type     int    `json:"type"`     // Error, warn
	Payload  string `json:"payload"`
}

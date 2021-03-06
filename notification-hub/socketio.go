package main

import (
	"encoding/json"
	"log"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/streadway/amqp"
)

var server *gosocketio.Server

func initSocketServer() {
	server = gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Channel %s created", c.Id())
		c.Join(MySpaceGeneralChannel)
		c.Join(BroadcastChannel)
	})

	server.On("identify", func(c *gosocketio.Channel, username string) string {
		log.Printf("Channel %s identified as %s", c.Id(), username)
		c.Join(username)
		return "OK"
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Channel %s disconnected", c.Id())
	})

}

func BroadcastTo(message amqp.Delivery) {
	log.Printf(" [x] %s", message.Exchange)
	log.Printf(" [x] %s", message.RoutingKey)
	log.Printf(" [x] %s", message.Body)

	var content RabbitMqMsg

	if err := json.Unmarshal(message.Body, &content); err != nil {
		failOnError(err, "Couldn't unmarshal the message")
	} else {
		log.Printf(" [x] %s", content.To)
		log.Printf(" [x] %s", content.Function)
		server.BroadcastTo(content.To, string(content.Function), content)
	}

}

func serverToUser(message RabbitMqMsg) {
	log.Printf(" [x] %s", message.To)
	log.Printf(" [x] %s", string(message.Function))

	server.BroadcastTo(message.To, string(message.Function), message)
}

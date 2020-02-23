package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
)

// To move to an independant library
var connection *amqp.Connection
var queue amqp.Queue
var channel *amqp.Channel

func initRabbitMq() {
	// Get the connection string from the environment variable
	url := os.Getenv("AMQP_URL")

	//If it doesn't exist, use the default connection string.

	if url == "" {
		//Don't do this in production, this is for testing purposes only.
		url = "amqp://system-notifier:password@localhost:5672"
	}

	var err error
	// Connect to the rabbitMQ instance
	connection, err = amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err = connection.Channel()
	failOnError(err, "Failed to open a channel")

	err = channel.ExchangeDeclare(
		"user-notification", // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	queue, err = channel.QueueDeclare(
		"wmn-internal", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")
}

func sendMessage(exchange string, useQueue bool, data RabbitMqMsg) {
	body, err := json.Marshal(data)
	failOnError(err, "The object couldn't be marshalled")

	var routing string = data.To
	if useQueue {
		routing = queue.Name
	}

	log.Printf("Sending %d %d %d %d %s to: %s | %s", data.Status, data.Priority, data.Type, data.Function, data.Status, routing, body)

	err = channel.Publish(
		exchange,               // exchange
		"spaces.init."+routing, // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType:  "application/json; charset=UTF-8",
			Body:         []byte(body),
			DeliveryMode: 1,
		})

	failOnError(err, "Failed to publish a message")
}

func constructNotification(ids string, client string, function Function, status int, priority int, _type int, description string) RabbitMqMsg {
	return RabbitMqMsg{
		ID:       ids,
		To:       client,
		Status:   status,
		Function: function,
		Priority: priority,
		Type:     _type,
		Payload:  description,
	}

}

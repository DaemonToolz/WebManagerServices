package main

import (
	"os"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// To move to an independant library
var connection *amqp.Connection
var queue amqp.Queue
var channel *amqp.Channel
var messages <-chan amqp.Delivery

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
	//defer connection.Close()

	channel, err = connection.Channel()
	failOnError(err, "Failed to open a channel")
	//defer channel.Close()

	err = channel.ExchangeDeclare(
		"system-users-notification", // name
		"topic",                     // type
		true,                        // durable
		false,                       // auto-deleted
		false,                       // internal
		false,                       // no-wait
		nil,                         // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	messages, err = channel.Consume(
		"wmn-system-users-notification", // queue
		uuid.New().String(),             // consumer
		false,                           // auto ack
		false,                           // exclusive
		false,                           // no local
		false,                           // no wait
		nil,                             // args
	)
	failOnError(err, "Failed to register a consumer")

}

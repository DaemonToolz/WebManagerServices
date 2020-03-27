package main

import (
	"os"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// To move to an independant library
var connection *amqp.Connection
var queue amqp.Queue
var channels []*amqp.Channel
var messages []<-chan amqp.Delivery

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
	channels = make([]*amqp.Channel, 2)
	channels[0], err = connection.Channel()
	failOnError(err, "Failed to open a channel")
	channels[1], err = connection.Channel()
	failOnError(err, "Failed to open a channel")
	//defer channel.Close()

	// duplicate
	err = channels[0].ExchangeDeclare(
		"user-notification", // name
		"topic",             // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)

	failOnError(err, "Failed to declare an exchange")

	err = channels[1].ExchangeDeclare(
		"account-notifications", // name
		"topic",                 // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)

	failOnError(err, "Failed to declare an exchange")

	messages = make([]<-chan amqp.Delivery, 2)
	messages[0], err = channels[0].Consume(
		"wmn-user-notification", // queue
		uuid.New().String(),     // consumer
		false,                   // auto ack
		false,                   // exclusive
		false,                   // no local
		false,                   // no wait
		nil,                     // args
	)

	messages[1], err = channels[1].Consume(
		"wmn-account-notification", // queue
		uuid.New().String(),        // consumer
		false,                      // auto ack
		false,                      // exclusive
		false,                      // no local
		false,                      // no wait
		nil,                        // args
	)

	failOnError(err, "Failed to register a consumer")

}

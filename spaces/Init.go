package main

import (
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
	// Git repos here
)

func main() {
	prepareLogs();
	initConfiguration();
	defer logFile.Close();
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
	defer connection.Close()

	channel, err = connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

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

	router := NewRouter()
	log.Fatal(http.ListenAndServe(appConfig.httpListenUri(), router))

}

package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// To move to an independant library
var connection *amqp.Connection
var queue amqp.Queue
var channel *amqp.Channel

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendMessage(exchange string, useQueue bool, data RabbitMqMsg) {
	body, err := json.Marshal(data)
	failOnError(err, "The object couldn't be marshalled")

	var routing string = data.To
	if useQueue {
		routing = queue.Name
	}

	err = channel.Publish(
		exchange, // exchange
		routing,  // routing key
		false,    // mandatory
		true,     // immediate
		amqp.Publishing{
			ContentType: "application/json; charset=UTF-8",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func constructMessage(client string, function string, status int, priority int, _type int, data interface{}) RabbitMqMsg {
	return RabbitMqMsg{
		To:       client,
		Status:   status,
		Priority: priority,
		Type:     _type,
	}
}

func constructNotification(client string, function string, status int, priority int, _type int) RabbitMqMsg {
	return RabbitMqMsg{
		To:       client,
		Function: function,
		Priority: priority,
		Type:     _type,
	}
}

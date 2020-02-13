package main

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

var connection amqp.Connection
var queue amqp.Queue
var channel amqp.Channel

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendMessage(data RabbitMqMsg) {
	body, err := json.Marshal(data)
	failOnError(err, "The object couldn't be marshalled")
	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json; charset=UTF-8",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

/*
func constructMessage(client string, status int, priority int, _type:int, data interface{}) RabbitMqMsg {
	message := RabbitMqMsg{
		To:client,
		Status:status,
		Priority:priority,
		Type:_type,

	}
	return message
}

func constructNotificationn(client string,  priority int, _type:int) RabbitMqMsg {
	message := RabbitMqMsg{
		To:client,
		Priority:priority,
		Type:_type,
	}
	return message;
}
*/

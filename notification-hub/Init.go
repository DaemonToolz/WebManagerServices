package main

import (
	"log"
	"net/http"
)

func main() {
	initConfiguration()
	prepareLogs()
	initRabbitMq()
	initSocketServer()

	go func() {
		for message := range messages {
			BroadcastTo(message)
			message.Ack(true)
		}
	}()

	defer logFile.Close()
	defer connection.Close()
	defer channel.Close()

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	log.Println("Serving at ", appConfig.httpListenUri(), "/socket.io/")
	log.Fatal(http.ListenAndServe(appConfig.httpListenUri(), serveMux))

}

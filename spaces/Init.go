package main

import (
	"log"
	"net/http"
	// Git repos here
)

func main() {

	prepareLogs()
	log.Println("Spaces service started")

	initConfiguration()
	log.Println("Configurations initialized")

	defer logFile.Close()
	defer connection.Close()
	defer channel.Close()

	initRabbitMq()
	log.Println("RabbitMQ initialized")
	router := NewRouter()
	initMiddleware(router)

	log.Fatal(http.ListenAndServe(appConfig.httpListenUri(), router))
	log.Println("Spaces service ended")
}

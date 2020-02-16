package main

import (
	"log"
	"net/http"
)

func main() {

	initRabbitMq()
	initSocketServer()

	go func() {
		for message := range messages {
			BroadcastTo(message)
		}
	}()

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	log.Println("Serving at localhost:20000...")
	log.Fatal(http.ListenAndServe(":20000", serveMux))

}

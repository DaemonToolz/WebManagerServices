package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"syscall"
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

	serveMux := http.NewServeMux()
	go serveMux.Handle("/socket.io/", server)

	log.Println("Serving at ", appConfig.httpListenUri(), "/socket.io/")
	go log.Fatal(http.ListenAndServe(appConfig.httpListenUri(), serveMux))

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, os.Kill)

	select {
	case <-sigChan:
		logFile.Close()
		connection.Close()
		channel.Close()
		os.Exit(0)
	}
}

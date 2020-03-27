package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	initConfiguration()
	prepareLogs()
	initRabbitMq()
	initSocketServer()

	for index := range messages {
		go func(sub_index int) {

			log.Println("Reading channel in index", sub_index)
			for message := range messages[sub_index] {
				log.Println("Message received by ", sub_index)
				BroadcastTo(message)
				message.Ack(true)
			}
		}(index)
	}
	serveMux := http.NewServeMux()
	go func() {
		serveMux.Handle("/socket.io/", server)
	}()

	go func() {
		log.Println("Serving at ", appConfig.httpListenUri(), "/socket.io/")
		log.Fatal(http.ListenAndServe(appConfig.httpListenUri(), serveMux))
	}()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, os.Kill)

	serverToUser(RabbitMqMsg{
		ID:       "0",
		Date:     time.Now(),
		Status:   STATUS_DONE,
		Function: NotifiationHubUpd,
		To:       string(BroadcastChannel),
		Priority: PRIORITY_HIGH,
		Type:     TYPE_SUCCESS,
		Payload:  "Hub online",
	})

	select {
	case <-sigChan:
		serverToUser(RabbitMqMsg{
			ID:       "0",
			Date:     time.Now(),
			Status:   STATUS_DONE,
			Function: NotifiationHubUpd,
			To:       string(BroadcastChannel),
			Priority: PRIORITY_CRITICAL,
			Type:     TYPE_ERROR,
			Payload:  "Hub shutting down",
		})

		logFile.Close()
		connection.Close()
		for index := range channels {
			channels[index].Close()
		}
		os.Exit(0)
	}
}

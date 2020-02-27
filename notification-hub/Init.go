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

	go func() {
		for message := range messages {
			BroadcastTo(message)
			message.Ack(true)
		}
	}()

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

	time.Sleep(5 * time.Second)
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
		channel.Close()
		os.Exit(0)
	}
}

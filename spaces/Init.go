package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"syscall"

	"github.com/google/uuid"
	// Git repos here
)

func main() {

	prepareLogs()
	log.Println("Spaces service started")

	initConfiguration()
	log.Println("Configurations initialized")

	initRabbitMq()
	log.Println("RabbitMQ initialized")
	router := NewRouter()
	initMiddleware(router)

	initRemoteProcedureCall()
	log.Println("Filewatch monitoring initialized")

	periodicCheck()
	log.Println("Watchers initialized")

	go func() {
		log.Fatal(http.ListenAndServe(appConfig.httpListenUri(), router))
	}()

	log.Println("HTTP Server online")
	sendMessage("user-notification", false, constructNotification(uuid.New().String(), MySpaceGeneralChannel, MySpaceNotify, STATUS_NEW, PRIORITY_CRITICAL, TYPE_INFO, "MySpace service online"))

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, os.Kill)
	select {
	case <-sigChan:
		sendMessage("user-notification", false, constructNotification(uuid.New().String(), MySpaceGeneralChannel, MySpaceNotify, STATUS_NEW, PRIORITY_CRITICAL, TYPE_ERROR, "MySpace service offline"))

		logFile.Close()
		connection.Close()
		channel.Close()
		globalTimer.Stop()
		os.Exit(0)
	}
}

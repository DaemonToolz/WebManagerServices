package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"syscall"
	// Git repos here
)

func main() {
	prepareLogs()
	log.Println("Logs Ready")

	initConfiguration()
	log.Println("Configurations initialized")

	initRabbitMq()
	log.Println("RabbitMQ Ready")
	Wrapper = ArangoWrapper{}
	log.Println("Preparing a new Arango Driver")

	Wrapper.initDriver("http://localhost:8529", "user-service", "password")
	defer Wrapper.Close()

	log.Println("Database ready")
	router := NewRouter()
	initMiddleware(router)

	go func() {
		log.Fatal(http.ListenAndServe(appConfig.httpListenUri(), router))
	}()

	log.Println("HTTP Server online")

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

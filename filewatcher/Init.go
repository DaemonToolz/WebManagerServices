package main

import (
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
	// Git repos here
)

func main() {

	log.Println(os.Args)
	if len(os.Args) != 1 {
		log.Println("Number of parameters not fitting. Shutting the watcher down")
		os.Exit(2)
	}

	user = os.Args[0]

	prepareLogs()
	log.Println("Spaces service started")

	initConfiguration()
	log.Println("Configurations initialized")

	defer logFile.Close()
	defer connection.Close()
	defer channel.Close()

	initRabbitMq()
	log.Println("RabbitMQ initialized")

	initFileWatcher(getPrivateFolders())
	//
	done := make(chan bool)
	defer watcher.Close()
	//
	go func() {
		for {

			select {
			// watch for events
			case event := <-watcher.Events:
				log.Printf("Event received %s", event)
				message := "The file " + event.Name + " : " + strconv.FormatInt(int64(event.Op), 10)
				sendMessage("user-notification", false, constructNotification(uuid.New().String(), user, "FileWatch", STATUS_NEW, PRIORITY_STD, TYPE_INFO, message))

				// watch for errors
			case err := <-watcher.Errors:
				failOnError(err, "The filewatcher detected an error")
				sendMessage("user-notification", false, constructNotification(uuid.New().String(), user, "FileWatch", STATUS_NEW, PRIORITY_STD, TYPE_INFO, err.Error()))
			}
		}
	}()

	<-done
	log.Println("Filewatcher service ended")
}

package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"

	"syscall"

	"github.com/google/uuid"
	// Git repos here
)

func main() {

	user = os.Args[1]

	prepareLogs()
	log.Println("Spaces service started")

	initConfiguration()
	log.Println("Configurations initialized")

	initRemoteProcedureCall()
	log.Println("RPC connected, connecting")

	initRabbitMq()
	log.Println("RabbitMQ initialized")

	initFileWatcher(getPrivateFolders())
	log.Println("Watcher initialized")

	sendMessage("user-notification", false, constructNotification(uuid.New().String(), user, FilewatchSysUpd, STATUS_DONE, PRIORITY_STD, TYPE_INFO, "Filewatch operational"))

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, os.Kill)

	go func() {
		for {

			select {
			// watch for events
			case event := <-watcher.Events:
				log.Printf("Event received %s", event)
				message := "The file " + event.Name + " : " + strconv.FormatInt(int64(event.Op), 10)
				sendMessage("user-notification", false, constructNotification(uuid.New().String(), user, FilewatchNotify, STATUS_NEW, PRIORITY_STD, TYPE_INFO, message))

				// watch for errors
			case err := <-watcher.Errors:
				failOnError(err, "The filewatcher detected an error")
				sendMessage("user-notification", false, constructNotification(uuid.New().String(), user, FilewatchNotify, STATUS_NEW, PRIORITY_STD, TYPE_INFO, err.Error()))
			}
		}
	}()

	select {
	case <-sigChan:
		log.Println("Filewatcher service ended")
		sendMessage("user-notification", false, constructNotification(uuid.New().String(), user, FilewatchSysUpd, STATUS_DONE, PRIORITY_STD, TYPE_WARN, "Filewatch stopped"))
		Unregister()
		log.Println("Filewatcher unregistered")
		watcher.Close()
		connection.Close()
		channel.Close()

		logFile.Close()
		os.Exit(0)
	}
}

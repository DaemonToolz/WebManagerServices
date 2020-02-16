package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	// Git repos here
)

func main() {

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/send", "payload", func(s socketio.Conn, mpayloadsg string) string {
		s.SetContext(payload)
		return payload
	})

	server.OnEvent("/", "disconnect", func(s socketio.Conn) string {
		s.Close()
		return ""
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	log.Println("Serving at localhost:20000...")
	log.Fatal(http.ListenAndServe(":20000", nil))

}

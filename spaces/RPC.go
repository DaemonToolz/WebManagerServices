package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type Args struct {
	PId  int
	User string
}

type RPCRegistry struct {
	Filewatches map[string]int `json:"filewatches"`
}

func (t *RPCRegistry) Register(args *Args, reply *int) error {
	t.Filewatches[args.User] = args.PId
	log.Println("Added filewatch user ", args.User, "/", args.PId)
	return nil
}

func IsRegistered(user string) bool {
	if _, ok := filewatchRegistry.Filewatches[user]; ok {
		return ok
	}
	return false
}

func (t *RPCRegistry) Unregister(args *Args, reply *int) error {
	delete(t.Filewatches, args.User)
	log.Println("Deleted filewatch user ", args.User, "/", args.PId)
	return nil
}

var filewatchRegistry *RPCRegistry

func initRemoteProcedureCall() {
	listener, err := net.Listen("tcp", appConfig.rpcListenUri())
	if err != nil {
		failOnError(err, "Couldn't initialize the RPC listener")
	}

	filewatchRegistry = &RPCRegistry{
		Filewatches: make(map[string]int),
	}
	rpc.Register(filewatchRegistry)
	rpc.HandleHTTP()

	log.Println("Opening RPC port for Filewatch system", listener.Addr().(*net.TCPAddr).Port)

	go http.Serve(listener, nil)
}

var globalTimer *time.Ticker

func periodicCheck() {
	globalTimer = time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-globalTimer.C:
				log.Println("Checking for unused watchers")
				initWatchers()
			}
		}
	}()
}

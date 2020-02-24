package main

import (
	"net/rpc"
	"os"
)

type Args struct {
	PId  int
	User string
}

var client *rpc.Client
var myself Args

func initRemoteProcedureCall() {
	openConnection()
	defer client.Close()
	myself = Args{
		PId:  os.Getpid(),
		User: user,
	}
	// Asynchronous call
	client.Go("RPCRegistry.Register", myself, nil, nil)
	// check errors, print, etc.

}

func openConnection() *rpc.Client {
	var err error
	client, err = rpc.DialHTTP("tcp", appConfig.rpcListenUri())
	if err != nil {
		failOnError(err, "couldn't connect to remote RPC server")
	}
	return client
}

func Unregister() {
	openConnection()
	defer client.Close()
	client.Go("RPCRegistry.Disconnect", myself, nil, nil)

}

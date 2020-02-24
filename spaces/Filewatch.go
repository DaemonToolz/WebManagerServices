package main

import (
	"log"
	"os"
)

func startFilewatch(user string) {
	attr := &os.ProcAttr{Dir: ".", Env: os.Environ(), Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
	process, err := os.StartProcess("filewatcher.exe", []string{user}, attr)
	if err == nil {

		// It is not clear from docs, but Realease actually detaches the process
		err = process.Release()
		if err != nil {
			failOnError(err, "An error occured when detaching the process")
		}

	} else {
		fmt.Println(err.Error())
	}

}

func initWatchers() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("A critical error occured : ", r)
		}
	}()

	for _, name := range getUsers() {
		if !IsRegistered(name) {
			startFilewatch(name)
		}
	}

}

package main

import (
	"log"
	"os"
	"syscall"
)

func startFilewatch(user string) {
	attr := &os.ProcAttr{Dir: ".", Env: os.Environ(), Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}

	process, err := os.StartProcess("filewatcher.exe", []string{user}, attr)
	if err == nil {
		err = process.Release()
		if err != nil {
			failOnError(err, "An error occured when detaching the process")
		}
	} else {
		failOnError(err, "An error has occured during the creation of the process")
	}

}

func initWatchers() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("A critical error occured : ", r)
		}
	}()

	for user, _ := range filewatchRegistry.Filewatches {
		exists, _ := exists(user)
		if !exists {
			clearWatchers(filewatchRegistry.Filewatches[user])
		}
	}

	for _, name := range getUsers() {
		if !IsRegistered(name) {
			startFilewatch(name)
		}
	}

}

func clearWatchers(pid int) {
	process, err := os.FindProcess(int(pid))
	if err != nil {
		log.Printf("Failed to find process: %s", err)
	} else {
		err := process.Signal(syscall.Signal(0))
		log.Printf("process.Signal on pid %d returned: %v\n", pid, err)
	}

}

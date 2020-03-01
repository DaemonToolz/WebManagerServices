package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func startFilewatch(user string) {
	copyExec(user)
	cmd := exec.Command("cmd.exe", user, "/d", ".", "/C", "start", "FILEWATCH_"+user+".exe", user)
	if err := cmd.Start(); err != nil {
		failOnError(err, "Couldn't start the desired process")
	}

}

func copyExec(user string) error {
	target := "FILEWATCH_" + user + ".exe"

	source, err := os.Open("filewatcher.exe")
	if err != nil {
		return err
	}
	defer source.Close()

	os.Remove(target)
	destination, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, 512)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err = destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

func initWatchers() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("A critical error occured : ", r)
		}
	}()

	for user, _ := range filewatchRegistry.Filewatches {
		exists, _ := exists(user)
		if !exists && !checkSpace(user).Created {
			clearWatchers(filewatchRegistry.Filewatches[user])
		}
		clearRegister(user)
	}

	for _, name := range getUsers() {
		if !IsRegistered(name) && checkSpace(name).Created {
			startFilewatch(name)
		}
	}

}

func clearWatchers(pid int) {
	process, err := os.FindProcess(int(pid))
	if err != nil {
		failOnError(err, "Process not found")
	} else {
		err := process.Signal(syscall.Signal(0))
		if err != nil {
			process.Kill()
		}
		process.Release()
		failOnError(err, "An error occured when closing the runner")
	}

}

func clearRegister(user string) {
	process, err := os.FindProcess(int(filewatchRegistry.Filewatches[user]))

	if err != nil {
		delete(filewatchRegistry.Filewatches, user)
		failOnError(err, "Error")
	} else {
		process.Release()
	}

}

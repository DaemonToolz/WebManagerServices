package main

import "os"

func startFilewatch(user string) {
	attr := &os.ProcAttr{Dir: ".", Env: os.Environ(), Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}}
	os.StartProcess("filewatcher.exe", []string{user}, attr)
}

func initWatchers() {
	//for names := range getUsers() {

	//}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	FolderRef string `json:"folderref`
}

var appConfig Config
var logFile os.File

var user string

const mySpacefolder = "myspace"
const privateFiles = "private"

func initConfiguration() {
	configFile, err := os.Open("./config/appConfig.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&appConfig)
}

func getPrivateFolders() string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", appConfig.FolderRef, privateFiles, user, mySpacefolder, mySpacefolder)
}

func constructHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func prepareLogs() {
	os.MkdirAll("./logs/", 0755)

	filename := fmt.Sprintf("./logs/%s.log", filepath.Base(os.Args[0]))
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	log.SetOutput(logFile)
}

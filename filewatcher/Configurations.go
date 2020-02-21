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

const mySpacefolder = "myspace"
const privateFiles = "private"
const sharedFiles = "shared"

func initConfiguration() {
	configFile, err := os.Open("./config/appConfig.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&appConfig)
}

func getPrivateFolders(username string) string {
	return fmt.Sprintf("%s/%s/%s/%s", appConfig.FolderRef, privateFiles, username, mySpacefolder)
}

func getSharedFolders() string {
	return fmt.Sprintf("%s/%s", appConfig.FolderRef, sharedFiles)
}

func getConfigurationFolder(username string) string {
	return fmt.Sprintf("%s/%s/%s", appConfig.FolderRef, privateFiles, username)
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

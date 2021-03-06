package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	RPCPort   int    `json:"rpcport"`
	FolderRef string `json:"folderref`
}

func (cfg *Config) httpListenUri() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

func (cfg *Config) rpcListenUri() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.RPCPort)
}

var appConfig Config
var logFile os.File

const mySpacefolder = "myspace"
const privateFiles = "private"
const sharedFiles = "shared"
const profileDataFolder = "profile_data"

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

func getUsersFolder() string {
	return fmt.Sprintf("%s/%s", appConfig.FolderRef, privateFiles)
}

func getExecPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
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

func checkSpace(space string) UserInitialization {
	configPath := fmt.Sprintf("%s/.%s.config.json", getConfigurationFolder(space), space)
	plan, _ := ioutil.ReadFile(configPath)

	var data UserInitialization
	err := json.Unmarshal(plan, &data)

	if err != nil {
		data = UserInitialization{
			UserId:     space,
			InitStatus: STATUS_ERROR,
			Created:    false,
		}
	}

	return data
}

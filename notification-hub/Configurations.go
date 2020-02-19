package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

var appConfig Config
var logFile os.File

func (cfg *Config) httpListenUri() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

func initConfiguration() {
	configFile, err := os.Open("./config/appConfig.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&appConfig)
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

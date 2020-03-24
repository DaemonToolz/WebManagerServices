package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func FindNetwork(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user string = vars["username"]

	data := GetNetwork(user)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf(err.Error())
		panic(err)
	}
}

func CheckUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	exists := userExists(user)

	if err := json.NewEncoder(w).Encode(exists); err != nil {
		log.Printf(err.Error())
		panic(err)
	}
}

func CheckOrCreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var userInfo UserInfo
	err := decoder.Decode(&userInfo)
	if err != nil {
		failOnError(err, "An error has occured")
	}

	log.Println("User ", userInfo, " has been received")

	created := struct {
		Success bool `json:"success"`
	}{false}

	exists := userExists(userInfo.Username).(struct {
		Exists bool `json:"exists"`
	})

	if !exists.Exists {
		log.Println("User ", userInfo, " does not exist")
		userInfo.CreatedAt = time.Now()
		go func() { CreateUser(userInfo) }()
		created.Success = true
	}

	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf(err.Error())
		panic(err)
	}
}

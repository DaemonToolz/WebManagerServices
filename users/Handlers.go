package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"fmt"
	"io"

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

		sendMessage(UsersRoot, AccountScope, ProfileCreated, AccountExchange, false,
			constructNotification("OK", userInfo.Username, UsersRegistered, int(STATUS_NEW), int(PRIORITY_MEDIUM), int(TYPE_SUCCESS), "Your profile has been created"))
	} else {
		sendMessage(UsersRoot, AccountScope, ProfileUpdated, AccountExchange, false,
			constructNotification("OK", userInfo.Username, UsersValidate, int(STATUS_NEW), int(PRIORITY_MEDIUM), int(TYPE_SUCCESS), "Your profile is up-to-date"))
	}

	if err := json.NewEncoder(w).Encode(created); err != nil {
		log.Printf(err.Error())
		panic(err)
	}
}

func GetProfilePicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var user string = vars["username"]

	client := http.Client{}
	// Utiliser un service de discovery
	requestedPicture, err := client.Get(fmt.Sprintf("http://localhost:10850/profile/picture/%s", user))
	log.Print(fmt.Sprintf("http://localhost:10850/profile/picture/%s", user))
	if err != nil {
		failOnError(err, "Couldn't fetch the profile picture for "+user)
		return
	}
	defer requestedPicture.Body.Close()

	log.Println("Attributes : ", fmt.Sprint(requestedPicture.ContentLength), " / ", requestedPicture.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", fmt.Sprint(requestedPicture.ContentLength))
	w.Header().Set("Content-Type", "image/png")
	if _, err = io.Copy(w, requestedPicture.Body); err != nil {
		log.Printf(err.Error())
		panic(err)
	}

}

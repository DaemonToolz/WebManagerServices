package main

import "log"

func userExists(username string) interface{} {
	log.Println("Checking if ", username, "exists")

	_, count := Wrapper.GetWhere(username)
	log.Println("Total count : ", count)
	return struct {
		Exists bool `json:"exists"`
	}{count > 0}
}

func CreateUser(data UserInfo) {
	log.Println("Creating user ", data)
	Wrapper.Create(data)
}

func GetNetwork(username string) []interface{} {
	log.Println("Checking network for ", username)

	data, count := Wrapper.GetConnected(username)
	log.Println("Total count : ", count)
	return data
}

package main

import "log"

func userExists(username string) interface{} {
	log.Println("Checking if ", username, "exists")

	_, count := Wrapper.GetWhere("_key", WhereEquals, UsersCollection, username)
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

	data, count := Wrapper.GetConnected(UsersGraph, UsersCollection, OutboundEdge, username, 1)
	log.Println("Total count : ", count)
	return data
}

func AddToNetwork(caller string, callee string, relation Relation) bool {
	log.Println("Adding ", callee, " to ", caller, " network")

	newRelation := RelationModel{
		From:     caller,
		To:       callee,
		Relation: string(relation),
	}

	success := Wrapper.AddRelation(UsersCollection, UsersEdge, newRelation)
	log.Println("Operation AddToNetwork : ", success)
	return success
}

func UpdateNetwork(caller string, callee string, relation Relation) bool {
	log.Println("Adding ", callee, " to ", caller, " network")

	newRelation := RelationModel{
		From:     caller,
		To:       callee,
		Relation: string(relation),
	}

	success := Wrapper.UpdateRelation(UsersCollection, UsersEdge, newRelation)
	log.Println("Operation AddToNetwork : ", success)
	return success
}

func RemoveFromNetwork(caller string, callee string, relation Relation) bool {
	log.Println("Adding ", callee, " to ", caller, " network")

	newRelation := RelationModel{
		From:     caller,
		To:       callee,
		Relation: string(relation),
	}

	success := Wrapper.RemoveRelation(UsersCollection, UsersEdge, newRelation)
	log.Println("Operation AddToNetwork : ", success)
	return success
}

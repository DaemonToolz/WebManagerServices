package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"FindNetwork",
		"GET",
		"/users/{username}/network",
		FindNetwork,
	},
	Route{
		"CreateUser",
		"POST",
		"/users/check",
		CheckOrCreateUser,
	},
}

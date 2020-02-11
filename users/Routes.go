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
		"LoadAccount",
		"GET",
		"/account/{id}",
		LoadAccount,
	},
	Route{
		"CreateAccount",
		"POST",
		"/account/create",
		CreateAccount,
	},
	Route{
		"login",
		"POST",
		"/account/login",
		Login,
	},
	Route{
		"validate",
		"POST",
		"/account/validate",
		ValidateAccount,
	},
}

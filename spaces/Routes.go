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
		"files",
		"GET",
		"/files/{space}",
		GetFiles,
	},
	Route{
		"file",
		"GET",
		"/files/{space}/{id}",
		GetFile,
	},
	Route{
		"create",
		"POST",
		"/spaces/init",
		CreateSpace,
	},
}

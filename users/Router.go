package main

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method, "OPTIONS").
			Path(route.Pattern).
			Name(route.Name).
			Handler(LoggerHandler(route.HandlerFunc))

	}

	return router
}

package main

import (
	"net/http"
	"github.com/codegangsta/negroni"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"validate",
		"/authorize/validate",
		negroni.New(
			negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
			negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				responseJSON(message, w, http.StatusOK)
		}))),
	},
}


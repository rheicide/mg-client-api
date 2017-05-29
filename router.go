package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Methods...).Path(route.Path).Handler(route.Handler)
	}

	allowedOrigins := handlers.AllowedOrigins([]string{
		"http://localhost:8080",
		"http://192.168.1.84:8080",
		"http://15.4.16.4:8080",
		"https://inbox.elarvee.xyz",
	})
	allowedMethods := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodOptions,
		http.MethodDelete,
	})
	allowedHeaders := handlers.AllowedHeaders([]string{
		"Authorization",
		"Content-Type",
	})

	return handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router)
}

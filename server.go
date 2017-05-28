package main

import (
	"log"

	"net/http"
)

func NewServer(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: NewRouter(),
	}
}

func StartServer(server *http.Server) {
	log.Printf("Starting server at %s", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

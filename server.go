package main

import (
	"log"
	"os"

	"net/http"

	r "gopkg.in/gorethink/gorethink.v3"
)

var (
	session *r.Session
)

func init() {
	var err error

	session, err = r.Connect(r.ConnectOpts{
		Address:  os.Getenv("R_ADDR"),
		Database: "mailgun",
	})

	if err != nil {
		log.Fatalln(err)
	}
}

func NewServer(addr string) *http.Server {
	return &http.Server{
		Addr: addr,
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

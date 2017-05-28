package main

import (
	"encoding/json"
	"net/http"

	"log"
	"time"

	"errors"
	"fmt"

	"strings"

	"crypto/rsa"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	r "gopkg.in/gorethink/gorethink.v3"
)

type HttpError struct {
	Err    error
	Status int
}

// implement error interface
func (e HttpError) Error() string {
	return e.Err.Error()
}

type Handler func(http.ResponseWriter, *http.Request) error

var (
	rsaPublicKey *rsa.PublicKey
)

func init() {
	publicKey := `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvXCk70W1gEOz699tXYv/NkKjiT9FV97o+pj+gNWzBpaoyv4S3QNC+I8pW0sVu5qNygtDJ72x1aDA
gWrOMYNg1OC8JiYvQLdEYcYpTy9m8RObM+Cpz/iHVGnEdPS8jxqJ27kTIBG1joQ2HyVbYDZfWUHIK1ks0pnXZvuSTUGD0/qge8hu1EaRQDuv/rA/y3XObNTi
Khcz8gGvCtDtdsvlDuEwfOgugGujHFpATlhLvfzhzbV5MznUhX89p+Lzf7j+XqWaoDaLScUgvzAo6vBs3pXfswWTMxYqv3SFkFqEmDLQNfx724n04GP1BMYU
rccXtUlC/6GN0b4Rro4ncAiArQIDAQAB
-----END PUBLIC KEY-----`

	var err error
	rsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		log.Fatalln(err)
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if err := validateToken(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	} else {
		err = h(w, r)
		if err != nil {
			var status int

			switch err := err.(type) {
			case HttpError:
				status = err.Status
			default:
				status = http.StatusInternalServerError
			}

			http.Error(w, err.Error(), status)
		}
	}

	log.Printf("%s\t%s\t%s", r.Method, r.RequestURI, time.Since(start))
}

func validateToken(r *http.Request) error {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return errors.New("Token is missing")
	} else {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	}

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return rsaPublicKey, nil
	})

	return err
}

func MailIndex(w http.ResponseWriter, _ *http.Request) error {
	res, err := r.Table("mails").
		Pluck("id", "from", "subject", "date").
		OrderBy(r.Desc("date")).
		Limit(10).
		Run(session)
	if err != nil {
		return HttpError{err, http.StatusInternalServerError}
	}

	var mails Mails
	if err = res.All(&mails); err != nil {
		return HttpError{err, http.StatusInternalServerError}
	}

	if err = json.NewEncoder(w).Encode(&mails); err != nil {
		return HttpError{err, http.StatusInternalServerError}
	}

	return nil
}

func MailById(w http.ResponseWriter, req *http.Request) error {
	id := mux.Vars(req)["id"]
	res, err := r.Table("mails").Get(id).Run(session)
	if err != nil {
		return HttpError{err, http.StatusInternalServerError}
	}

	var mail Mail
	if err = res.One(&mail); err != nil {
		return HttpError{err, http.StatusNotFound}
	}

	if err = json.NewEncoder(w).Encode(&mail); err != nil {
		return HttpError{err, http.StatusInternalServerError}
	}

	return nil
}

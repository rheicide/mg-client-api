package main

import (
	"encoding/json"
	"log"
	"time"
)

type Mail struct {
	Id        string    `json:"id" gorethink:"id,omitempty"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Recipient string    `json:"recipient"`
	Subject   string    `json:"subject"`
	BodyPlain string    `json:"body_plain,omitempty" gorethink:"body-plain"`
	BodyHtml  string    `json:"body_html,omitempty" gorethink:"body-html"`
	Date      time.Time `json:"date"`
	Read      bool      `json:"read"`
	Starred   bool      `json:"starred"`
}

// custom marshaller to convert date into my timezone
func (mail *Mail) MarshalJSON() ([]byte, error) {
	zone, _ := mail.Date.Zone()
	if zone == "+07:00" {
		return json.Marshal(*mail)
	}

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Fatalln(err)
	}

	type Alias Mail

	return json.Marshal(&struct {
		*Alias
		Date time.Time `json:"date"`
	}{
		Alias: (*Alias)(mail),
		Date:  mail.Date.In(loc),
	})
}

type Mails []Mail

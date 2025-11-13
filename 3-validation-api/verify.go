package main

import (
	"github.com/jordan-wright/email"
	"net/http"
	"net/smtp"
)

type Email struct {
	text string
}

func NewEmailHandler(Mux *http.ServeMux) {
	email := Email{}
	Mux.HandleFunc("/send", email.Send)
	Mux.HandleFunc("/verify/{hash}", email.Verify)

}
func (m Email) Send(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)

}

func (m Email) Verify(w http.ResponseWriter, r *http.Request) {

}

func CreateAndSendEmail(text string) {
	e := email.NewEmail()
	e.From = "Daniel Kosovskiy <broiler747tb@gmail.com>"
	e.To = []string{"churkabes420@gmail.com"}
	e.Subject = "Verification Email"
	e.Text = []byte(text)
	e.Send("smtp.gmail.com:587", smtp.PlainAuth("Daniel", "broiler747tb@gmail.com", "xxx", "smtp.gmail.com"))
}

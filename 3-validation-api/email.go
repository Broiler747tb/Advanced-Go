package main

import "github.com/jordan-wright/email"

func CreateEmailConfig(address string) *email.Email {
	e := email.NewEmail()
	e.From = "Daniel Kosovskiy <broiler747tb@gmail.com>"
	e.To = []string{address}
	e.Subject = "Verification email test"
	return e
}

type Email struct {
	Email string `json:"Email"`
}

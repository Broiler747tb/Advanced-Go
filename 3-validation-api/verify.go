package main

import (
	"encoding/json"
	"fmt"
	"io"
	rand2 "math/rand/v2"
	"net/http"
	"net/smtp"
	"os"
)

type User struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
	Hash     string `json:"Hash"`
}

func NewEmailHandler(Mux *http.ServeMux) {
	email := Email{}
	Mux.HandleFunc("/send", email.Send)
	Mux.HandleFunc("/verify/{hash}", email.Verify)

}
func (m Email) Send(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received user, sending verification request...")
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body: ", err)
		_, err = w.Write([]byte("Error reading request body"))
		if err != nil {
			fmt.Println(err)
		}
	}
	data := Email{}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println("Error unmarshalling json: ", err)
	}
	if data.Email == "" {
		fmt.Println("Error, empty email address")
		_, err = w.Write([]byte("Error, empty email address"))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		hash := MakeUser(data)
		address := fmt.Sprint("http://localhost:8081/verify/", hash)
		User{Email: data.Email}.CreateAndSendEmail(address)
	}

}

func (m Email) Verify(w http.ResponseWriter, r *http.Request) {
	hash := r.PathValue("hash")
	file, err := os.ReadFile("user.json")
	if err != nil {
		fmt.Println(err)
	}
	user := User{}
	err = json.Unmarshal(file, &user)
	if err != nil {
		fmt.Println(err)
	}
	if user.Hash != hash {
		fmt.Println("ERROR: USER HASH DOES NOT MATCH")
	} else {
		fmt.Println("USER AUTHENTIFIED")
	}
	w.WriteHeader(200)
}

func (u User) CreateAndSendEmail(text string) {
	mail := CreateEmailConfig(u.Address)
	mail.Text = []byte(text)
	err := mail.Send("smtp.gmail.com:587", smtp.PlainAuth("Daniel", "broiler747tb@gmail.com", "xxx", "smtp.gmail.com"))
	if err != nil {
		fmt.Println(err)
	}
}

func GenerateHash() string {
	var hash string
	for i := 0; i < 31; i++ {
		hash = hash + string(rune(rand2.Int32N(25)+97))
	}
	return hash
}

func MakeUser(mail Email) string {
	hash := GenerateHash()
	usr := User{Email: mail.Email, Hash: hash}
	bytes, err := json.Marshal(usr)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("user.json", bytes, 0775)
	if err != nil {
		fmt.Println(err)
	}
	return hash
}

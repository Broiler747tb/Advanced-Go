package main

import (
	"fmt"
	"net/http"
)

func main() {
	Mux := http.NewServeMux()
	NewEmailHandler(Mux)
	server := http.Server{
		Addr:    ":8081",
		Handler: Mux,
	}
	fmt.Println("Listening and serving on port: 8081")
	server.ListenAndServe()
}

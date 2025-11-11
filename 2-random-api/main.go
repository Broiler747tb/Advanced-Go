package main

import (
	rand2 "math/rand/v2"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/random", rand)
	http.ListenAndServe(":8081", nil)

}

func rand(w http.ResponseWriter, req *http.Request) {
	random := rand2.Int32N(6) + 1
	ran := strconv.Itoa(int(random))
	w.Write([]byte(ran))
}

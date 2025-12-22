package res

import (
	"encoding/json"
	"log"
	"net/http"
)

func Json(w http.ResponseWriter, data any, statuscode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Failed to give .json response", err)
	}
}

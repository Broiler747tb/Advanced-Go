package req

import (
	"encoding/json"
	"fmt"
	"io"
)

func Decode[T any](body io.ReadCloser) (T, error) {
	var payload T
	err := json.NewDecoder(body).Decode(&payload)
	if err != nil {
		fmt.Println("Error here", err)
		return payload, err
	}
	return payload, nil
}

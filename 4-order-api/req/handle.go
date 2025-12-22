package req

import (
	"net/http"
	"order-api/res"
)

func HandleBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {

		res.Json(w, err.Error(), 402)
		return nil, err
	}
	err = IsValid(body)
	if err != nil {
		res.Json(w, err.Error(), 402)
		return nil, err
	}
	return &body, nil
}

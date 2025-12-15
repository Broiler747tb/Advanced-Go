package Product

import "github.com/lib/pq"

type LinkCreateRequest struct {
	Name string `json:"name" validate:"required"`
}

type LinkUpdateRequest struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description"`
	Hash        string         `json:"hash"`
	Images      pq.StringArray `json:"images"`
}

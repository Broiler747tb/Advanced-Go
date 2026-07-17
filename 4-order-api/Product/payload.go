package Product

type LinkCreateRequest struct {
	Name string `json:"name" validate:"required"`
}

type LinkUpdateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Hash        string `json:"hash" gorm:"uniqueIndex"`
	Images      string `json:"images"`
}

package Order

type OrderCreateRequest struct {
	ProductIds []uint `json:"productIds" validate:"required,min=1"`
}

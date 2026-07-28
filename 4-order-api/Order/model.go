package Order

import (
	"order-api/Product"
	"order-api/User"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID   uint              `json:"userId"`
	User     User.User         `json:"-"`
	Products []Product.Product `json:"products" gorm:"many2many:order_products"`
}

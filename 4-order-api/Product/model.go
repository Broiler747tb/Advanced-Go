package Product

import (
	"gorm.io/gorm"
	rand2 "math/rand/v2"
)

type Product struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Hash        string `json:"hash" gorm:"uniqueIndex"`
	ImageLink   string `json:"imageLink"`
}

func NewProduct(name string) *Product {
	prod := &Product{Name: name}
	prod.GenerateHash()
	return prod
}

func (product *Product) GenerateHash() {
	var hash string
	for i := 0; i < 31; i++ {
		hash = hash + string(rune(rand2.Int32N(25)+97))
	}
	product.Hash = hash
}

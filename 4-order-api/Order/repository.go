package Order

import (
	"order-api/Db"
)

type OrderRepository struct {
	Database *Db.Db
}

func NewOrderRepository(database *Db.Db) *OrderRepository {
	return &OrderRepository{Database: database}
}

func (repo *OrderRepository) Create(order *Order) (*Order, error) {
	result := repo.Database.DB.Create(order)
	if result.Error != nil {
		return nil, result.Error
	}
	return order, nil
}

func (repo *OrderRepository) GetById(id uint) (*Order, error) {
	var order Order
	result := repo.Database.DB.Preload("Products").First(&order, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (repo *OrderRepository) GetByUserId(userId uint) ([]Order, error) {
	var orders []Order
	result := repo.Database.DB.Preload("Products").Find(&orders, "user_id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

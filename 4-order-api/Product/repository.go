package Product

import (
	"fmt"
	"order-api/Db"
)

type ProductRepository struct {
	Database *Db.Db
}

func NewProductRepository(database *Db.Db) *ProductRepository {
	return &ProductRepository{Database: database}
}

func (repo *ProductRepository) Create(product *Product) (*Product, error) {
	result := repo.Database.DB.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (repo *ProductRepository) GetByHash(hash string) (*Product, error) {
	if repo.Database == nil {
		return nil, fmt.Errorf("repo.Database is nil")
	}
	if repo.Database.DB == nil {
		return nil, fmt.Errorf("repo.Database.DB is nil")
	}

	var product Product
	result := repo.Database.DB.First(&product, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (repo *ProductRepository) GetById(id uint) (*Product, error) {
	var product Product
	result := repo.Database.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (repo *ProductRepository) Update(product *Product) (*Product, error) {
	result := repo.Database.DB.Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (repo *ProductRepository) Delete(id uint) error {
	result := repo.Database.DB.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

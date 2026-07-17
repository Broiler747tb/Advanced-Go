package User

import "order-api/Db"

type UserRepository struct {
	Database *Db.Db
}

func NewUserRepository(database *Db.Db) *UserRepository {
	return &UserRepository{Database: database}
}

func (repo *UserRepository) Create(product *User) (*User, error) {
	result := repo.Database.DB.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	var product User
	result := repo.Database.DB.First(&product, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (repo *UserRepository) FindByPhone(phone int) (*User, error) {
	var product User
	result := repo.Database.DB.First(&product, "phone = ?", phone)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

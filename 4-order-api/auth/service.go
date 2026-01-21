package auth

import (
	"errors"
	"order-api/User"
)

type AuthService struct {
	UserRepository User.UserRepository
}

func NewAuthService(UserRpository User.UserRepository) *AuthService {
	return &AuthService{UserRepository: UserRpository}
}

func (service *AuthService) Register(email, password, name string, phone int) (string, error) {
	userExists, _ := service.UserRepository.FindByPhone(phone)
	if userExists != nil {
		return "", errors.New("user already exists")
	}

	user := &User.User{
		Phone:    phone,
		Name:     name,
		Email:    email,
		Password: "",
	}
	_, err := service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

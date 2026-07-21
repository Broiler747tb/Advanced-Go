package auth

import (
	"errors"
	"fmt"

	"order-api/User"
	"order-api/pkg/jwt"
)

type UserRepo interface {
	FindByPhone(phone int) (*User.User, error)
	Create(user *User.User) (*User.User, error)
}

type AuthService struct {
	UserRepository UserRepo
	Sessions       *SessionStore
	Secret         string
}

func NewAuthService(userRepository UserRepo, secret string) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
		Sessions:       NewSessionStore(),
		Secret:         secret,
	}
}

func (service *AuthService) SendCode(phone int) (string, error) {
	code := GenerateCode()
	sessionId := service.Sessions.Create(phone, code)

	fmt.Printf("[SMS] phone=%d code=%d session=%s\n", phone, code, sessionId)
	return sessionId, nil
}

func (service *AuthService) VerifyCode(sessionId string, code int) (string, error) {
	phone, ok := service.Sessions.Verify(sessionId, code)
	if !ok {
		return "", errors.New("invalid session or code")
	}

	if existing, _ := service.UserRepository.FindByPhone(phone); existing == nil {
		if _, err := service.UserRepository.Create(&User.User{Phone: phone}); err != nil {
			return "", err
		}
	}

	return jwt.NewJWT(service.Secret).Create(jwt.JWTData{Phone: phone})
}

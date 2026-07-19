package auth

import (
	"errors"
	"fmt"

	"order-api/User"
	"order-api/pkg/jwt"
)

// UserRepo is the slice of the user repository the auth service depends on.
// Depending on an interface (rather than the concrete *User.UserRepository)
// lets tests inject an in-memory fake, so the auth flow can be exercised
// without a real database.
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

// SendCode generates a confirmation code, stores a session, "sends" the SMS
// (logged to the console for the lesson) and returns the session id.
func (service *AuthService) SendCode(phone int) (string, error) {
	code := GenerateCode()
	sessionId := service.Sessions.Create(phone, code)
	// Real app: hand off to an SMS provider here. Lesson: log it so it's readable.
	fmt.Printf("[SMS] phone=%d code=%d session=%s\n", phone, code, sessionId)
	return sessionId, nil
}

// VerifyCode validates the code for the session, find-or-creates the user by
// phone, and returns a signed JWT.
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

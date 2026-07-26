package jwt

import (
	"errors"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Phone int
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{Secret: secret}
}

func (j *JWT) Create(data JWTData) (string, error) {
	t := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, jwtv5.MapClaims{
		"phone": data.Phone,
	})
	return t.SignedString([]byte(j.Secret))
}

func (j *JWT) Parse(tokenString string) (bool, *JWTData) {
	t, err := jwtv5.Parse(tokenString, func(t *jwtv5.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.Secret), nil
	})
	if err != nil || !t.Valid {
		return false, nil
	}
	claims, ok := t.Claims.(jwtv5.MapClaims)
	if !ok {
		return false, nil
	}
	phone, _ := claims["phone"].(float64)
	return true, &JWTData{Phone: int(phone)}
}

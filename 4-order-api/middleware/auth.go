package middleware

import (
	"context"
	"net/http"
	"strings"

	"order-api/pkg/jwt"
)

type key string

const ContextPhoneKey key = "ContextPhoneKey"

func IsAuthed(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		valid, data := jwt.NewJWT(secret).Parse(token)
		if !valid {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextPhoneKey, data.Phone)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PhoneFromContext(ctx context.Context) (int, bool) {
	phone, ok := ctx.Value(ContextPhoneKey).(int)
	return phone, ok
}

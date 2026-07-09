package middleware

import (
	"context"
	"net/http"
	"strings"

	"order-api/pkg/jwt"
)

type key string

const ContextPhoneKey key = "ContextPhoneKey"

// IsAuthed wraps a handler, requiring a valid "Authorization: Bearer <jwt>"
// header. On success the caller's phone is stored on the request context.
func IsAuthed(next http.Handler, secret string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		valid, data := jwt.NewJWT(secret).Parse(token)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ContextPhoneKey, data.Phone)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

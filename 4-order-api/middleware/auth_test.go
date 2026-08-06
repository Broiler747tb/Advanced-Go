package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"order-api/pkg/jwt"
)

const secret = "test-secret"

func newToken(t *testing.T, signedWith string, phone int) string {
	t.Helper()
	token, err := jwt.NewJWT(signedWith).Create(jwt.JWTData{Phone: phone})
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestIsAuthedAllowsValidToken(t *testing.T) {
	phone := 0
	called := false
	handler := IsAuthed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		phone, _ = PhoneFromContext(r.Context())
	}), secret)

	request := httptest.NewRequest(http.MethodPost, "/link", nil)
	request.Header.Set("Authorization", "Bearer "+newToken(t, secret, 89990009900))
	writer := httptest.NewRecorder()

	handler.ServeHTTP(writer, request)

	if !called {
		t.Fatal("handler was not called")
	}
	if writer.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", writer.Code, http.StatusOK)
	}
	if phone != 89990009900 {
		t.Errorf("phone = %d, want %d", phone, 89990009900)
	}
}

func TestIsAuthedRejectsBadTokens(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{"no header", ""},
		{"no bearer prefix", newToken(t, secret, 89990009900)},
		{"garbage token", "Bearer abc.def.ghi"},
		{"foreign secret", "Bearer " + newToken(t, "other-secret", 89990009900)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := IsAuthed(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("handler should not be reached")
			}), secret)

			request := httptest.NewRequest(http.MethodPost, "/link", nil)
			if test.header != "" {
				request.Header.Set("Authorization", test.header)
			}
			writer := httptest.NewRecorder()

			handler.ServeHTTP(writer, request)

			if writer.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want %d", writer.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestPhoneFromContextWithoutMiddleware(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/link", nil)
	if _, ok := PhoneFromContext(request.Context()); ok {
		t.Error("expected no phone in a plain request context")
	}
}

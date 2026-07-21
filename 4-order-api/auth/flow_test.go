package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-api/Configs"
	"order-api/User"
	"order-api/middleware"
	"order-api/pkg/jwt"
)

type fakeUserRepo struct {
	byPhone map[int]*User.User
	created int
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{byPhone: map[int]*User.User{}}
}

func (f *fakeUserRepo) FindByPhone(phone int) (*User.User, error) {
	if u, ok := f.byPhone[phone]; ok {
		return u, nil
	}

	return nil, errors.New("record not found")
}

func (f *fakeUserRepo) Create(user *User.User) (*User.User, error) {
	f.byPhone[user.Phone] = user
	f.created++
	return user, nil
}

func buildRouter(secret string, repo UserRepo) (*http.ServeMux, *AuthService) {
	router := http.NewServeMux()
	service := NewAuthService(repo, secret)

	NewAuthHandler(router, AuthHandlerDeps{
		Config:      &Configs.Config{Auth: Configs.AuthConfig{Secret: secret}},
		AuthService: service,
	})

	protected := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		phone := r.Context().Value(middleware.ContextPhoneKey)
		fmt.Fprintf(w, "hello %v", phone)
	})
	router.Handle("GET /me", middleware.IsAuthed(protected, secret))

	return router, service
}

func postJSON(t *testing.T, router http.Handler, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestSMSAuthFlow_HappyPath(t *testing.T) {
	const secret = "test-secret"
	const phone = "89990009900"
	repo := newFakeUserRepo()
	router, service := buildRouter(secret, repo)

	rec := postJSON(t, router, "/auth/send-code", map[string]string{"phone": phone})
	if rec.Code != http.StatusOK {
		t.Fatalf("send-code: got status %d, want 200; body=%s", rec.Code, rec.Body)
	}
	var sent SendCodeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &sent); err != nil {
		t.Fatalf("send-code: decode response: %v", err)
	}
	if sent.SessionId == "" {
		t.Fatal("send-code: expected a non-empty sessionId")
	}

	service.Sessions.mu.Lock()
	code := service.Sessions.sessions[sent.SessionId].Code
	service.Sessions.mu.Unlock()
	if code < 1000 || code > 9999 {
		t.Fatalf("expected a 4-digit code, got %d", code)
	}

	rec = postJSON(t, router, "/auth/verify-code", map[string]any{
		"sessionId": sent.SessionId,
		"code":      code,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("verify-code: got status %d, want 200; body=%s", rec.Code, rec.Body)
	}
	var verified VerifyCodeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &verified); err != nil {
		t.Fatalf("verify-code: decode response: %v", err)
	}
	if verified.Token == "" {
		t.Fatal("verify-code: expected a non-empty token")
	}

	valid, data := jwt.NewJWT(secret).Parse(verified.Token)
	if !valid {
		t.Fatal("verify-code: returned token did not validate as a real JWT")
	}
	if data.Phone != 89990009900 {
		t.Fatalf("token phone = %d, want 89990009900", data.Phone)
	}

	if repo.created != 1 {
		t.Fatalf("expected exactly 1 user created, got %d", repo.created)
	}

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+verified.Token)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /me with valid token: got %d, want 200", rec.Code)
	}
	if rec.Body.String() != "hello 89990009900" {
		t.Fatalf("GET /me body = %q, want %q", rec.Body.String(), "hello 89990009900")
	}
}

func TestVerifyCode_WrongCodeIsRejected(t *testing.T) {
	const secret = "test-secret"
	router, service := buildRouter(secret, newFakeUserRepo())

	rec := postJSON(t, router, "/auth/send-code", map[string]string{"phone": "89990009900"})
	var sent SendCodeResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &sent)

	service.Sessions.mu.Lock()
	realCode := service.Sessions.sessions[sent.SessionId].Code
	service.Sessions.mu.Unlock()

	wrong := realCode + 1
	rec = postJSON(t, router, "/auth/verify-code", map[string]any{
		"sessionId": sent.SessionId,
		"code":      wrong,
	})
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("wrong code: got status %d, want 401", rec.Code)
	}
}

func TestProtectedEndpoint_RejectsMissingAndBadToken(t *testing.T) {
	const secret = "test-secret"
	router, _ := buildRouter(secret, newFakeUserRepo())

	cases := []struct {
		name   string
		header string
	}{
		{"no header", ""},
		{"not bearer", "Basic abc"},
		{"garbage token", "Bearer not-a-real-jwt"},
		{"wrong secret", "Bearer " + mustSign(t, "other-secret", 123)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != http.StatusUnauthorized {
				t.Fatalf("%s: got status %d, want 401", tc.name, rec.Code)
			}
		})
	}
}

func mustSign(t *testing.T, secret string, phone int) string {
	t.Helper()
	token, err := jwt.NewJWT(secret).Create(jwt.JWTData{Phone: phone})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return token
}

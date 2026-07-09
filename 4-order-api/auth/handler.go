package auth

import (
	"net/http"
	"strconv"

	"order-api/Configs"
	"order-api/req"
	"order-api/res"
)

type AuthHandlerDeps struct {
	*Configs.Config
	*AuthService
}

type AuthHandler struct {
	*Configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/send-code", handler.SendCode())
	router.HandleFunc("POST /auth/verify-code", handler.VerifyCode())
}

func (handler *AuthHandler) SendCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[SendCodeRequest](w, r)
		if err != nil {
			return
		}
		phone, err := strconv.Atoi(body.Phone) // phone arrives as a JSON string
		if err != nil {
			http.Error(w, "invalid phone", http.StatusBadRequest)
			return
		}
		// qualify .AuthService: handler.SendCode is this method, not the service's
		sessionId, err := handler.AuthService.SendCode(phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, SendCodeResponse{SessionId: sessionId}, http.StatusOK)
	}
}

func (handler *AuthHandler) VerifyCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[VerifyCodeRequest](w, r)
		if err != nil {
			return
		}
		token, err := handler.AuthService.VerifyCode(body.SessionId, body.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		res.Json(w, VerifyCodeResponse{Token: token}, http.StatusOK)
	}
}

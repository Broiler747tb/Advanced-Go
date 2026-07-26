package auth

type SendCodeRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type SendCodeResponse struct {
	SessionId string `json:"sessionId"`
}

type VerifyCodeRequest struct {
	SessionId string `json:"sessionId" validate:"required"`
	Code      int    `json:"code" validate:"required"`
}

type VerifyCodeResponse struct {
	Token string `json:"token"`
}

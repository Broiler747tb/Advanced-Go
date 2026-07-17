package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	LoginRequest
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Phone    int    `json:"phone"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegisterResponce struct {
	Token string `json:"token"`
}

package controller

// RegisterRequest é o body de POST /auth/register.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest é o body de POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse é a resposta de POST /auth/login.
type LoginResponse struct {
	Token string `json:"token"`
}

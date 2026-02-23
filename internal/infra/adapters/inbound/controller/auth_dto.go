package controller

import "fivestars/internal/domain"

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

func ToDomainRegister(req RegisterRequest) (*domain.UserRegistration, error) {
	userResistration := &domain.UserRegistration{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	err := userResistration.Validate()
	if err != nil {
		return nil, err
	}

	return userResistration, nil
}

func ToDomainLogin(req LoginRequest) (*domain.UserCredentials, error) {
	userCredentials := &domain.UserCredentials{
		Email:    req.Email,
		Password: req.Password,
	}

	err := userCredentials.Validate()
	if err != nil {
		return nil, err
	}

	return userCredentials, nil
}

func ToLoginResponse(authResult domain.AuthenticationResult) (*LoginResponse, error) {
	loginResponse := &LoginResponse{
		Token: authResult.Token,
	}

	return loginResponse, nil
}

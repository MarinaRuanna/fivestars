package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/application/usecases"
)

const authBodyMaxBytes int64 = 32 << 10 // 32KB

type AuthHandler struct {
	registerUserUseCase usecases.RegisterUserUseCase
	loginUserUseCase    usecases.LoginUserUseCase
}

// NewAuthHandler cria um AuthHandler.
func NewAuthHandler(
	registerUserUseCase usecases.RegisterUserUseCase,
	loginUserUseCase usecases.LoginUserUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUserUseCase: registerUserUseCase,
		loginUserUseCase:    loginUserUseCase,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var req RegisterRequest
	if err := decodeStrictJSONBody(w, r, &req, authBodyMaxBytes); err != nil {
		return err
	}

	registration, err := ToDomainRegister(req)
	if err != nil {
		return err
	}

	output, err := h.registerUserUseCase.Execute(r.Context(), *registration)
	if err != nil {
		return err
	}

	registrationDTO, err := ToLoginResponse(*output)
	if err != nil {
		return err
	}

	// 5. FORMAT HTTP RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(LoginResponse{Token: registrationDTO.Token}); err != nil {
		return err
	}

	return nil
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest
	if err := decodeStrictJSONBody(w, r, &req, authBodyMaxBytes); err != nil {
		return err
	}

	credentials, err := ToDomainLogin(req)
	if err != nil {
		return err
	}

	output, err := h.loginUserUseCase.Execute(r.Context(), *credentials)
	if err != nil {
		return err
	}

	LoginResponseDTO, err := ToLoginResponse(*output)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(LoginResponse{Token: LoginResponseDTO.Token}); err != nil {
		return err
	}

	return nil
}

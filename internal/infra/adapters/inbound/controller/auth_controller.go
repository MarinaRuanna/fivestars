package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/application/usecases"
	"fivestars/internal/domain/customerror"
)

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

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, customerror.NewValidationError("body JSON inválido"))
		return
	}

	registration, err := ToDomainRegister(req)
	if err != nil {
		respondError(w, err)
		return
	}

	output, err := h.registerUserUseCase.Execute(r.Context(), *registration)
	if err != nil {
		respondError(w, err)
		return
	}

	registrationDTO, err := ToLoginResponse(*output)
	if err != nil {
		respondError(w, err)
		return
	}

	// 5. FORMAT HTTP RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(LoginResponse{Token: registrationDTO.Token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, customerror.NewValidationError("body JSON inválido"))
		return
	}

	credentials, err := ToDomainLogin(req)
	if err != nil {
		respondError(w, err)
		return
	}

	output, err := h.loginUserUseCase.Execute(r.Context(), *credentials)
	if err != nil {
		respondError(w, err)
		return
	}

	LoginResponseDTO, err := ToLoginResponse(*output)
	if err != nil {
		respondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(LoginResponse{Token: LoginResponseDTO.Token})
}

func respondError(w http.ResponseWriter, err error) {
	code := customerror.StatusCodeFromError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

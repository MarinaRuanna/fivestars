package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/application/usecases"
	"fivestars/internal/domain/customerror"
)

// AuthHandler trataPOST /auth/register e POST /auth/login.
// ⭐ REFATORADO: Agora apenas ORQUESTRA HTTP + delegação para usecases
type AuthHandler struct {
	registerUserUseCase *usecases.RegisterUserUseCase
	loginUserUseCase    *usecases.LoginUserUseCase
}

// NewAuthHandler cria um AuthHandler.
func NewAuthHandler(
	registerUserUseCase *usecases.RegisterUserUseCase,
	loginUserUseCase *usecases.LoginUserUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUserUseCase: registerUserUseCase,
		loginUserUseCase:    loginUserUseCase,
	}
}

// Register trata POST /auth/register.
// ⭐ REFATORADO: Handler apenas parseia HTTP + chama usecase
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. PARSE HTTP
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, customerror.NewValidationError("body JSON inválido"))
		return
	}

	// 2. DELEGAR PARA USECASE (toda a lógica de negócio)
	output, err := h.registerUserUseCase.Execute(r.Context(), usecases.RegisterUserInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})

	// 3. HANDLE ERRO
	if err != nil {
		respondError(w, err)
		return
	}

	// 4. FORMAT HTTP RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(LoginResponse{Token: output.Token})
}

// Login trata POST /auth/login.
// ⭐ REFATORADO: Handler apenas parseia HTTP + chama usecase
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. PARSE HTTP
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, customerror.NewValidationError("body JSON inválido"))
		return
	}

	// 2. DELEGAR PARA USECASE (toda a lógica de negócio)
	output, err := h.loginUserUseCase.Execute(r.Context(), usecases.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	})

	// 3. HANDLE ERRO
	if err != nil {
		respondError(w, err)
		return
	}

	// 4. FORMAT HTTP RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(LoginResponse{Token: output.Token})
}

// respondError escreve erro no formato JSON e status apropriado usando customerror.
func respondError(w http.ResponseWriter, err error) {
	code := customerror.StatusCodeFromError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

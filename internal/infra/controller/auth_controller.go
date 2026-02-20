package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/auth"
)

// AuthHandler trata POST /auth/register e POST /auth/login.
type AuthHandler struct {
	userRepo domain.UserRepository
	jwtSecret string
}

// NewAuthHandler cria um AuthHandler.
func NewAuthHandler(userRepo domain.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, jwtSecret: jwtSecret}
}

// Register trata POST /auth/register.
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
	if req.Email == "" || req.Password == "" || req.Name == "" {
		respondError(w, customerror.NewValidationError("email, password e name são obrigatórios"))
		return
	}
	if len(req.Password) < 6 {
		respondError(w, customerror.NewValidationError("password deve ter no mínimo 6 caracteres"))
		return
	}
	existing, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		respondError(w, err)
		return
	}
	if existing != nil {
		respondError(w, customerror.NewConflictError("email já cadastrado"))
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	u := &domain.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
		Level:        1,
	}
	if err := u.Validate(); err != nil {
		respondError(w, err)
		return
	}
	if err := h.userRepo.Create(r.Context(), u); err != nil {
		respondError(w, err)
		return
	}
	token, err := auth.NewToken(u.ID, h.jwtSecret, 0)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// Login trata POST /auth/login.
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
	if req.Email == "" || req.Password == "" {
		respondError(w, customerror.NewValidationError("email e password são obrigatórios"))
		return
	}
	u, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		respondError(w, err)
		return
	}
	if u == nil || !auth.CheckPassword(u.PasswordHash, req.Password) {
		respondError(w, customerror.NewUnauthorizedError("credenciais inválidas"))
		return
	}
	token, err := auth.NewToken(u.ID, h.jwtSecret, 0)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// respondError escreve erro no formato JSON e status apropriado usando customerror.
func respondError(w http.ResponseWriter, err error) {
	code := customerror.StatusCodeFromError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/domain"
	"fivestars/internal/infra/auth"
)

// UserHandler trata GET /users/me (requer autenticação).
type UserHandler struct {
	userRepo domain.UserRepository
}

// NewUserHandler cria um UserHandler.
func NewUserHandler(userRepo domain.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// Me retorna o perfil do usuário logado. Deve ser chamado após o middleware de auth.
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		return
	}
	u, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if u == nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(UserFromDomain(u))
}

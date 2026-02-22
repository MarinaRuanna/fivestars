package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/application"
	"fivestars/internal/infra/auth"
)

// UserHandler trata GET /users/me (requer autenticação).
// ⭐ REFATORADO: Agora apenas ORQUESTRA HTTP + delegação para usecase
type UserHandler struct {
	getUserUseCase *application.GetUserUseCase
}

// NewUserHandler cria um UserHandler.
func NewUserHandler(getUserUseCase *application.GetUserUseCase) *UserHandler {
	return &UserHandler{getUserUseCase: getUserUseCase}
}

// Me retorna o perfil do usuário logado. Deve ser chamado após o middleware de auth.
// ⭐ REFATORADO: Handler apenas extrai userID + chama usecase
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. EXTRACT FROM CONTEXT (set by auth middleware)
	userID := auth.UserIDFromContext(r.Context())

	// 2. DELEGAR PARA USECASE
	output, err := h.getUserUseCase.Execute(r.Context(), application.GetUserInput{UserID: userID})
	if err != nil {
		// Format error response properly
		w.Header().Set("Content-Type", "application/json")
		if err.Error() == "user not authenticated" {
			w.WriteHeader(http.StatusUnauthorized)
		} else if err.Error() == "user not found" {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 3. FORMAT HTTP RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(output)
}

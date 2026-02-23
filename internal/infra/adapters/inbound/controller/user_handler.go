package controller

import (
	"encoding/json"
	"fivestars/internal/application/usecases"
	"fivestars/internal/infra/auth"
	"net/http"
)

type UserHandler struct {
	getUserUseCase usecases.GetUserUseCase
}

func NewUserHandler(getUserUseCase usecases.GetUserUseCase) *UserHandler {
	return &UserHandler{getUserUseCase: getUserUseCase}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := auth.UserIDFromContext(r.Context())

	output, err := h.getUserUseCase.Execute(r.Context(), userID)
	if err != nil {
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

	userDTO := UserFromDomain(output)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(userDTO)
}

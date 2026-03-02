package controller

import (
	"encoding/json"
	"fivestars/internal/application/usecases"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/auth"
	"net/http"
)

type UserHandler struct {
	getUserUseCase usecases.GetUserUseCase
}

func NewUserHandler(getUserUseCase usecases.GetUserUseCase) *UserHandler {
	return &UserHandler{getUserUseCase: getUserUseCase}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		return customerror.NewUnauthorizedError("user not authenticated")
	}

	output, err := h.getUserUseCase.Execute(r.Context(), userID)
	if err != nil {
		return err
	}

	userDTO := UserFromDomain(output)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userDTO); err != nil {
		return err
	}

	return nil
}

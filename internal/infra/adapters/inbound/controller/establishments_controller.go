package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/application/usecases"
)

// EstablishmentsHandler handles GET /establishments.
// ⭐ REFATORADO: Agora apenas ORQUESTRA HTTP + delegação para usecase
type EstablishmentsHandler struct {
	listEstablishmentsUseCase *usecases.ListEstablishmentsUseCase
}

// NewEstablishmentsHandler returns a new EstablishmentsHandler.
func NewEstablishmentsHandler(
	listEstablishmentsUseCase *usecases.ListEstablishmentsUseCase,
) *EstablishmentsHandler {
	return &EstablishmentsHandler{listEstablishmentsUseCase: listEstablishmentsUseCase}
}

// List returns all establishments as JSON.
// ⭐ REFATORADO: Handler apenas chama usecase + formata resposta
func (h *EstablishmentsHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. DELEGAR PARA USECASE
	items, err := h.listEstablishmentsUseCase.Execute(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 2. FORMAT HTTP RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"items": items})
}

package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/domain"
)

// EstablishmentsHandler handles GET /establishments.
type EstablishmentsHandler struct {
	repo domain.EstablishmentRepository
}

// NewEstablishmentsHandler returns a new EstablishmentsHandler.
func NewEstablishmentsHandler(repo domain.EstablishmentRepository) *EstablishmentsHandler {
	return &EstablishmentsHandler{repo: repo}
}

// List returns all establishments as JSON usando DTO de resposta (proteção dos dados).
func (h *EstablishmentsHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	ctx := r.Context()
	list, err := h.repo.List(ctx)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	items := make([]EstablishmentResponse, 0, len(list))
	for _, e := range list {
		items = append(items, FromDomain(e))
	}
	resp := EstablishmentListResponse{Items: items}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

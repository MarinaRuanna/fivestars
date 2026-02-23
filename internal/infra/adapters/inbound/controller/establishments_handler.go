package controller

import (
	"encoding/json"
	"fivestars/internal/application/usecases"
	"net/http"
)

type EstablishmentsHandler struct {
	estabService usecases.ListEstablishmentsUseCase
}

func NewEstablishmentsHandler(
	estabService usecases.ListEstablishmentsUseCase,
) *EstablishmentsHandler {
	return &EstablishmentsHandler{estabService: estabService}
}

func (c *EstablishmentsHandler) ListEstablishments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	items, err := c.estabService.Execute(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	listDTO := FromDomainList(items)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"items": listDTO})
}

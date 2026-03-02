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

func (c *EstablishmentsHandler) ListEstablishments(w http.ResponseWriter, r *http.Request) error {
	items, err := c.estabService.Execute(r.Context())
	if err != nil {
		return err
	}

	listDTO := FromDomainList(items)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"items": listDTO}); err != nil {
		return err
	}

	return nil
}

package controller

import (
	"encoding/json"
	"net/http"

	"fivestars/internal/application"
)

// EstablishmentsController handles GET /establishments.
// ⭐ REFATORADO: Agora apenas ORQUESTRA HTTP + delegação para usecase
type EstablishmentsController struct {
	estabService application.EstablishmentService
}

// NewEstablishmentsController returns a new EstablishmentsController.
func NewEstablishmentsController(
	estabService application.EstablishmentService,
) *EstablishmentsController {
	return &EstablishmentsController{estabService: estabService}
}

// List returns all establishments as JSON.
// ⭐ REFATORADO: Handler apenas chama usecase + formata resposta
func (c *EstablishmentsController) ListEstablishments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. DELEGAR PARA USECASE
	items, err := c.estabService.ListEstablishments(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 2. FORMAT HTTP RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"items": items})
}

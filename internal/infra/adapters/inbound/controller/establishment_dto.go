package controller

import (
	"fivestars/internal/domain"
	"time"
)

// EstablishmentResponse é o contrato da API para um estabelecimento.
// Expõe apenas os campos desejados na resposta HTTP (proteção dos dados internos).
// QRCode não é exposto na API pública.
type EstablishmentResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	Category  string   `json:"category"`
	Address   string   `json:"address,omitempty"`
	Lat       *float64 `json:"lat,omitempty"`
	Lng       *float64 `json:"lng,omitempty"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// EstablishmentListResponse é o contrato da API para a lista de estabelecimentos.
type EstablishmentListResponse struct {
	Items []EstablishmentResponse `json:"items"`
}

// FromDomain converte uma entidade de domínio em DTO de resposta da API.
func FromDomain(e *domain.Establishment) EstablishmentResponse {
	return EstablishmentResponse{
		ID:        e.ID,
		Name:      e.Name,
		Slug:      e.Slug,
		Category:  e.Category,
		Address:   e.Address,
		Lat:       e.Lat,
		Lng:       e.Lng,
		CreatedAt: e.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: e.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

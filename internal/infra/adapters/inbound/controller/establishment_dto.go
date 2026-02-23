package controller

import (
	"fivestars/internal/domain"
	"time"
)

type EstablishmentResponse struct {
	ID        string   `json:"establishment_id"`
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	Category  string   `json:"category"`
	Address   string   `json:"address,omitempty"`
	Lat       *float64 `json:"lat,omitempty"`
	Lng       *float64 `json:"lng,omitempty"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type EstablishmentListResponse struct {
	Items []EstablishmentResponse `json:"items"`
}

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

func FromDomainList(estabs []domain.Establishment) []EstablishmentResponse {
	responses := make([]EstablishmentResponse, len(estabs))
	for i, estab := range estabs {
		responses[i] = FromDomain(&estab)
	}
	return responses
}

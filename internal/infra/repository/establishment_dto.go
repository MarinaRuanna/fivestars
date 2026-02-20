package repository

import (
	"time"

	"fivestars/internal/domain"
)

// EstablishmentRow representa uma linha da tabela establishments.
// Isola o formato do banco do domínio; o repositório faz a conversão Row → Domain.
type EstablishmentRow struct {
	ID        string
	Name      string
	Slug      string
	Category  string
	Address   string
	Lat       *float64
	Lng       *float64
	QRCode    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToDomain converte o DTO de persistência em entidade de domínio.
func (r *EstablishmentRow) ToDomain() *domain.Establishment {
	return &domain.Establishment{
		ID:        r.ID,
		Name:      r.Name,
		Slug:      r.Slug,
		Category:  r.Category,
		Address:   r.Address,
		Lat:       r.Lat,
		Lng:       r.Lng,
		QRCode:    r.QRCode,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

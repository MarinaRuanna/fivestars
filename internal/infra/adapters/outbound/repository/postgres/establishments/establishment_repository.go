package establishments

import (
	"context"

	"fivestars/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EstablishmentRepository handles persistence for establishments.
type establishmentRepository struct {
	pool *pgxpool.Pool
}

// NewEstablishmentRepository returns a new EstablishmentRepository.
func NewEstablishmentRepository(pool *pgxpool.Pool) domain.EstablishmentRepository {
	return &establishmentRepository{pool: pool}
}

// List returns all establishments (no pagination for Phase 1).
// Usa DTO de persistência (EstablishmentRow) e converte para domínio.
func (r *establishmentRepository) List(ctx context.Context) ([]*domain.Establishment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, slug, category, address, lat, lng, qr_code, created_at, updated_at
		FROM establishments
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*domain.Establishment
	for rows.Next() {
		var estabDTO EstablishmentDTO
		err := rows.Scan(
			&estabDTO.ID, &estabDTO.Name, &estabDTO.Slug, &estabDTO.Category, &estabDTO.Address,
			&estabDTO.Lat, &estabDTO.Lng, &estabDTO.QRCode, &estabDTO.CreatedAt, &estabDTO.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		estab, err := estabDTO.ToDomain()
		if err != nil {
			return nil, err
		}

		list = append(list, estab)
	}
	return list, rows.Err()
}

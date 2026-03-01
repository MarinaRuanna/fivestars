package establishments

import (
	"context"

	"fivestars/internal/domain"
	"fivestars/internal/infra/adapters/outbound/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type establishmentRepository struct {
	pool *pgxpool.Pool
}

func NewEstablishmentRepository(pool *pgxpool.Pool) domain.EstablishmentRepository {
	return &establishmentRepository{pool: pool}
}

func (r *establishmentRepository) List(ctx context.Context) ([]domain.Establishment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, slug, category, address, lat, lng, qr_code, created_at, updated_at
		FROM establishments
		ORDER BY name
	`)
	if err != nil {
		return nil, postgres.MapError(err, "establishment")
	}
	defer rows.Close()

	var list []domain.Establishment
	for rows.Next() {
		var estabDTO EstablishmentDTO
		err := rows.Scan(
			&estabDTO.ID, &estabDTO.Name, &estabDTO.Slug, &estabDTO.Category, &estabDTO.Address,
			&estabDTO.Lat, &estabDTO.Lng, &estabDTO.QRCode, &estabDTO.CreatedAt, &estabDTO.UpdatedAt,
		)
		if err != nil {
			return nil, postgres.MapError(err, "establishment")
		}
		estab, err := estabDTO.ToDomain()
		if err != nil {
			return nil, err
		}

		list = append(list, *estab)
	}
	if err := rows.Err(); err != nil {
		return nil, postgres.MapError(err, "establishment")
	}
	return list, nil
}

func (r *establishmentRepository) GetByID(ctx context.Context, establishmentID string) (*domain.Establishment, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, name, slug, category, address, lat, lng, qr_code, created_at, updated_at
		FROM establishments
		WHERE id = $1
	`, establishmentID)

	var estabDTO EstablishmentDTO
	if err := row.Scan(
		&estabDTO.ID, &estabDTO.Name, &estabDTO.Slug, &estabDTO.Category, &estabDTO.Address,
		&estabDTO.Lat, &estabDTO.Lng, &estabDTO.QRCode, &estabDTO.CreatedAt, &estabDTO.UpdatedAt,
	); err != nil {
		if postgres.IsNoRows(err) {
			return nil, nil
		}
		return nil, postgres.MapError(err, "establishment")
	}

	estab, err := estabDTO.ToDomain()
	if err != nil {
		return nil, err
	}
	return estab, nil
}

func (r *establishmentRepository) DistanceTo(ctx context.Context, id string, lat, lng float64) (float64, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT ST_DistanceSphere(ST_SetSRID(ST_MakePoint(lng, lat), 4326)::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) as meters
		FROM establishments
		WHERE id = $3
	`, lng, lat, id)

	var meters float64
	if err := row.Scan(&meters); err != nil {
		return 0, postgres.MapError(err, "establishment")
	}
	return meters, nil
}

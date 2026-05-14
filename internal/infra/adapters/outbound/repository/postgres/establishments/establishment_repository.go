package establishments

import (
	"context"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/adapters/outbound/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type establishmentRepository struct {
	pool *pgxpool.Pool
}

func NewEstablishmentRepository(pool *pgxpool.Pool) domain.EstablishmentRepository {
	return &establishmentRepository{pool: pool}
}

func (r *establishmentRepository) Create(ctx context.Context, establishment *domain.Establishment) error {
	dto, err := FromDomain(establishment)
	if err != nil {
		return err
	}

	err = r.pool.QueryRow(ctx, `
		INSERT INTO establishments (owner_id, name, slug, category, address, lat, lng, qr_code)
		VALUES (NULLIF($1, '')::uuid, $2, $3, $4, NULLIF($5, ''), $6, $7, NULLIF($8, ''))
		RETURNING id, created_at, updated_at
	`, dto.OwnerID, dto.Name, dto.Slug, dto.Category, dto.Address, dto.Lat, dto.Lng, dto.QRCode).Scan(&establishment.ID, &establishment.CreatedAt, &establishment.UpdatedAt)
	if err != nil {
		return postgres.MapError(err, "establishment")
	}

	return nil
}

func (r *establishmentRepository) ClaimOwnership(ctx context.Context, establishmentID, ownerID, claimQRCode string) (*domain.Establishment, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, postgres.MapError(err, "establishment")
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `
		SELECT id, COALESCE(owner_id::text, ''), name, slug, category,
		       COALESCE(address, '') as address,
		       lat, lng,
		       COALESCE(qr_code, '') as qr_code,
		       created_at, updated_at
		FROM establishments
		WHERE id = $1
		FOR UPDATE
	`, establishmentID)

	var dto EstablishmentDTO
	if err := row.Scan(
		&dto.ID, &dto.OwnerID, &dto.Name, &dto.Slug, &dto.Category, &dto.Address,
		&dto.Lat, &dto.Lng, &dto.QRCode, &dto.CreatedAt, &dto.UpdatedAt,
	); err != nil {
		if postgres.IsNoRows(err) {
			return nil, nil
		}
		return nil, postgres.MapError(err, "establishment")
	}

	if dto.OwnerID != "" && dto.OwnerID != ownerID {
		return nil, customerror.NewConflictError("establishment already claimed")
	}
	if dto.OwnerID == ownerID {
		if dto.QRCode == "" || dto.QRCode != claimQRCode {
			return nil, customerror.NewForbiddenError("invalid claim qr_code")
		}
		establishment, err := dto.ToDomain()
		if err != nil {
			return nil, err
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, postgres.MapError(err, "establishment")
		}
		return establishment, nil
	}
	if dto.QRCode == "" || dto.QRCode != claimQRCode {
		return nil, customerror.NewForbiddenError("invalid claim qr_code")
	}

	row = tx.QueryRow(ctx, `
		UPDATE establishments
		SET owner_id = NULLIF($2, '')::uuid,
		    updated_at = NOW()
		WHERE id = $1 AND owner_id IS NULL
		RETURNING id, COALESCE(owner_id::text, ''), name, slug, category,
		          COALESCE(address, '') as address,
		          lat, lng,
		          COALESCE(qr_code, '') as qr_code,
		          created_at, updated_at
	`, establishmentID, ownerID)

	if err := row.Scan(
		&dto.ID, &dto.OwnerID, &dto.Name, &dto.Slug, &dto.Category, &dto.Address,
		&dto.Lat, &dto.Lng, &dto.QRCode, &dto.CreatedAt, &dto.UpdatedAt,
	); err != nil {
		if postgres.IsNoRows(err) {
			return nil, customerror.NewConflictError("establishment already claimed")
		}
		return nil, postgres.MapError(err, "establishment")
	}

	establishment, err := dto.ToDomain()
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, postgres.MapError(err, "establishment")
	}

	return establishment, nil
}

func (r *establishmentRepository) List(ctx context.Context) ([]domain.Establishment, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, COALESCE(owner_id::text, ''), name, slug, category,
		       COALESCE(address, '') as address,
		       lat, lng,
		       COALESCE(qr_code, '') as qr_code,
		       created_at, updated_at
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
			&estabDTO.ID, &estabDTO.OwnerID, &estabDTO.Name, &estabDTO.Slug, &estabDTO.Category, &estabDTO.Address,
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
		SELECT id, COALESCE(owner_id::text, ''), name, slug, category,
		       COALESCE(address, '') as address,
		       lat, lng,
		       COALESCE(qr_code, '') as qr_code,
		       created_at, updated_at
		FROM establishments
		WHERE id = $1
	`, establishmentID)

	var estabDTO EstablishmentDTO
	if err := row.Scan(
		&estabDTO.ID, &estabDTO.OwnerID, &estabDTO.Name, &estabDTO.Slug, &estabDTO.Category, &estabDTO.Address,
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

func (r *establishmentRepository) GetStats(ctx context.Context, establishmentID string) (*domain.EstablishmentStats, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT
			COALESCE(review_stats.average_rating, 0),
			COALESCE(review_stats.total_reviews, 0),
			COALESCE(like_stats.total_likes, 0),
			COALESCE(highlight_stats.highlighted_review_count, 0)
		FROM establishments e
		LEFT JOIN (
			SELECT establishment_id, AVG(rating)::float8 AS average_rating, COUNT(*) AS total_reviews
			FROM reviews
			GROUP BY establishment_id
		) review_stats ON review_stats.establishment_id = e.id
		LEFT JOIN (
			SELECT r.establishment_id, COUNT(*) AS total_likes
			FROM reviews r
			JOIN review_likes rl ON rl.review_id = r.id
			GROUP BY r.establishment_id
		) like_stats ON like_stats.establishment_id = e.id
		LEFT JOIN (
			SELECT establishment_id, COUNT(*) AS highlighted_review_count
			FROM highlights
			GROUP BY establishment_id
		) highlight_stats ON highlight_stats.establishment_id = e.id
		WHERE e.id = $1
	`, establishmentID)

	var dto EstablishmentStatsDTO
	if err := row.Scan(&dto.AverageRating, &dto.TotalReviews, &dto.TotalLikes, &dto.HighlightedReviewCount); err != nil {
		if postgres.IsNoRows(err) {
			return nil, nil
		}
		return nil, postgres.MapError(err, "establishment")
	}

	return dto.ToDomain()
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

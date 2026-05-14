package highlights

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/adapters/outbound/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type highlightRepository struct {
	pool *pgxpool.Pool
}

func NewHighlightRepository(pool *pgxpool.Pool) domain.HighlightRepository {
	return &highlightRepository{pool: pool}
}

func (r *highlightRepository) Create(ctx context.Context, highlight *domain.Highlight) error {
	dto, err := FromDomain(highlight)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO highlights (establishment_id, review_id, created_by_user_id, created_at)
		VALUES ($1, $2, $3, $4)
	`, dto.EstablishmentID, dto.ReviewID, dto.CreatedByUserID, dto.CreatedAt)
	if err != nil {
		return postgres.MapError(fmt.Errorf("insert highlight: %w", err), "highlight")
	}

	return nil
}

func (r *highlightRepository) Delete(ctx context.Context, establishmentID, reviewID string) error {
	cmd, err := r.pool.Exec(ctx, `
		DELETE FROM highlights
		WHERE establishment_id = $1 AND review_id = $2
	`, establishmentID, reviewID)
	if err != nil {
		return postgres.MapError(fmt.Errorf("delete highlight: %w", err), "highlight")
	}
	if cmd.RowsAffected() == 0 {
		return customerror.NewNotFoundError("highlight not found")
	}

	return nil
}

func (r *highlightRepository) Exists(ctx context.Context, establishmentID, reviewID string) (bool, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM highlights
			WHERE establishment_id = $1 AND review_id = $2
		)
	`, establishmentID, reviewID)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, postgres.MapError(err, "highlight")
	}

	return exists, nil
}

func (r *highlightRepository) CountByEstablishment(ctx context.Context, establishmentID string) (int, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM highlights
		WHERE establishment_id = $1
	`, establishmentID)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, postgres.MapError(err, "highlight")
	}

	return count, nil
}

func (r *highlightRepository) ListByEstablishment(ctx context.Context, establishmentID string) ([]domain.Highlight, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT establishment_id, review_id, created_by_user_id, created_at
		FROM highlights
		WHERE establishment_id = $1
		ORDER BY created_at DESC
	`, establishmentID)
	if err != nil {
		return nil, postgres.MapError(err, "highlight")
	}
	defer rows.Close()

	var list []domain.Highlight
	for rows.Next() {
		var dto HighlightDTO
		if err := rows.Scan(&dto.EstablishmentID, &dto.ReviewID, &dto.CreatedByUserID, &dto.CreatedAt); err != nil {
			return nil, postgres.MapError(err, "highlight")
		}

		highlight, err := dto.ToDomain()
		if err != nil {
			return nil, err
		}

		list = append(list, *highlight)
	}
	if err := rows.Err(); err != nil {
		return nil, postgres.MapError(err, "highlight")
	}

	return list, nil
}

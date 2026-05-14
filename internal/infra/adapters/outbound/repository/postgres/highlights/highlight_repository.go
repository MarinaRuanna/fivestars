package highlights

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/adapters/outbound/repository/postgres"

	"github.com/jackc/pgx/v5"
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

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return postgres.MapError(fmt.Errorf("begin highlight transaction: %w", err), "highlight")
	}
	defer tx.Rollback(ctx)

	var lock int
	if err := tx.QueryRow(ctx, `
		SELECT 1
		FROM establishments
		WHERE id = $1
		FOR UPDATE
	`, dto.EstablishmentID).Scan(&lock); err != nil {
		if postgres.IsNoRows(err) {
			return customerror.NewNotFoundError("establishment not found")
		}
		return postgres.MapError(fmt.Errorf("lock establishment for highlight: %w", err), "highlight")
	}

	cmd, err := tx.Exec(ctx, `
		INSERT INTO highlights (establishment_id, review_id, created_by_user_id, created_at)
		SELECT $1, $2, $3, $4
		WHERE (
			SELECT COUNT(*)
			FROM highlights
			WHERE establishment_id = $1
		) < $5
		AND NOT EXISTS (
			SELECT 1
			FROM highlights
			WHERE establishment_id = $1 AND review_id = $2
		)
	`, dto.EstablishmentID, dto.ReviewID, dto.CreatedByUserID, dto.CreatedAt, domain.MaxHighlightsPerEstablishment)
	if err != nil {
		return postgres.MapError(fmt.Errorf("insert highlight: %w", err), "highlight")
	}
	if cmd.RowsAffected() == 0 {
		exists, err := r.existsTx(ctx, tx, dto.EstablishmentID, dto.ReviewID)
		if err != nil {
			return postgres.MapError(fmt.Errorf("check highlight conflict: %w", err), "highlight")
		}
		if exists {
			return customerror.NewConflictError("review already highlighted")
		}

		count, err := r.countByEstablishmentTx(ctx, tx, dto.EstablishmentID)
		if err != nil {
			return postgres.MapError(fmt.Errorf("check highlight limit: %w", err), "highlight")
		}
		if count >= domain.MaxHighlightsPerEstablishment {
			return customerror.NewConflictError("highlight limit reached")
		}

		return customerror.NewConflictError("highlight could not be created")
	}

	if err := tx.Commit(ctx); err != nil {
		return postgres.MapError(fmt.Errorf("commit highlight transaction: %w", err), "highlight")
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

func (r *highlightRepository) existsTx(ctx context.Context, tx pgx.Tx, establishmentID, reviewID string) (bool, error) {
	row := tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM highlights
			WHERE establishment_id = $1 AND review_id = $2
		)
	`, establishmentID, reviewID)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (r *highlightRepository) countByEstablishmentTx(ctx context.Context, tx pgx.Tx, establishmentID string) (int, error) {
	row := tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM highlights
		WHERE establishment_id = $1
	`, establishmentID)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
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

package reviewlikes

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"
	"fivestars/internal/infra/adapters/outbound/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type reviewLikeRepository struct {
	pool *pgxpool.Pool
}

func NewReviewLikeRepository(pool *pgxpool.Pool) domain.ReviewLikeRepository {
	return &reviewLikeRepository{pool: pool}
}

func (r *reviewLikeRepository) Create(ctx context.Context, like *domain.ReviewLike) error {
	dto, err := FromDomain(like)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO review_likes (user_id, review_id, created_at)
		VALUES ($1, $2, $3)
	`, dto.UserID, dto.ReviewID, dto.CreatedAt)
	if err != nil {
		return postgres.MapError(fmt.Errorf("insert review like: %w", err), "review like")
	}
	return nil
}

func (r *reviewLikeRepository) Delete(ctx context.Context, userID, reviewID string) error {
	cmd, err := r.pool.Exec(ctx, `
		DELETE FROM review_likes WHERE user_id = $1 AND review_id = $2
	`, userID, reviewID)
	if err != nil {
		return postgres.MapError(fmt.Errorf("delete review like: %w", err), "review like")
	}
	if cmd.RowsAffected() == 0 {
		return customerror.NewNotFoundError("review like not found")
	}
	return nil
}

package reviews

import (
	"context"
	"fmt"

	"fivestars/internal/domain"
	"fivestars/internal/infra/adapters/outbound/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type reviewRepository struct {
	pool *pgxpool.Pool
}

func NewReviewRepository(pool *pgxpool.Pool) domain.ReviewRepository {
	return &reviewRepository{pool: pool}
}

func (r *reviewRepository) Create(ctx context.Context, review *domain.Review) error {
	dto, err := FromDomain(review)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, `
		INSERT INTO reviews (id, checkin_id, user_id, establishment_id, rating, text, created_at, updated_at)
		VALUES (COALESCE(NULLIF($1, ''), uuid_generate_v4()), $2, $3, $4, $5, $6, $7, $8)
	`, dto.ID, dto.CheckinID, dto.UserID, dto.EstablishmentID, dto.Rating, dto.Text, dto.CreatedAt, dto.UpdatedAt)
	if err != nil {
		return postgres.MapError(fmt.Errorf("insert review: %w", err), "review")
	}

	return nil
}

func (r *reviewRepository) GetByID(ctx context.Context, reviewID string) (*domain.Review, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT r.id, r.checkin_id, r.user_id, r.establishment_id, r.rating, r.text, r.created_at, r.updated_at,
		       COALESCE(l.like_count, 0) as like_count
		FROM reviews r
		LEFT JOIN (
			SELECT review_id, COUNT(*) as like_count
			FROM review_likes
			GROUP BY review_id
		) l ON l.review_id = r.id
		WHERE r.id = $1
	`, reviewID)

	var dto ReviewDTO
	if err := row.Scan(&dto.ID, &dto.CheckinID, &dto.UserID, &dto.EstablishmentID, &dto.Rating, &dto.Text, &dto.CreatedAt, &dto.UpdatedAt, &dto.LikeCount); err != nil {
		if postgres.IsNoRows(err) {
			return nil, nil
		}
		return nil, postgres.MapError(err, "review")
	}

	return dto.ToDomain()
}

func (r *reviewRepository) GetByCheckinID(ctx context.Context, checkinID string) (*domain.Review, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, checkin_id, user_id, establishment_id, rating, text, created_at, updated_at
		FROM reviews
		WHERE checkin_id = $1
		LIMIT 1
	`, checkinID)

	var dto ReviewDTO
	if err := row.Scan(&dto.ID, &dto.CheckinID, &dto.UserID, &dto.EstablishmentID, &dto.Rating, &dto.Text, &dto.CreatedAt, &dto.UpdatedAt); err != nil {
		if postgres.IsNoRows(err) {
			return nil, nil
		}
		return nil, postgres.MapError(err, "review")
	}

	return dto.ToDomain()
}

func (r *reviewRepository) ListByEstablishment(ctx context.Context, establishmentID string, options domain.ReviewListOptions) ([]domain.Review, error) {
	query := `
		SELECT r.id, r.checkin_id, r.user_id, r.establishment_id, r.rating, r.text, r.created_at, r.updated_at,
		       COALESCE(l.like_count, 0) as like_count
		FROM reviews r
		LEFT JOIN (
			SELECT review_id, COUNT(*) as like_count
			FROM review_likes
			GROUP BY review_id
		) l ON l.review_id = r.id
		WHERE r.establishment_id = $1
	`
	args := []any{establishmentID}
	argPos := 2

	if options.MinRating != nil {
		query += fmt.Sprintf(" AND r.rating >= $%d", argPos)
		args = append(args, *options.MinRating)
		argPos++
	}
	if options.MaxRating != nil {
		query += fmt.Sprintf(" AND r.rating <= $%d", argPos)
		args = append(args, *options.MaxRating)
		argPos++
	}

	query += " ORDER BY r.created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, options.Limit, options.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, postgres.MapError(err, "review")
	}
	defer rows.Close()

	var list []domain.Review
	for rows.Next() {
		var dto ReviewDTO
		if err := rows.Scan(&dto.ID, &dto.CheckinID, &dto.UserID, &dto.EstablishmentID, &dto.Rating, &dto.Text, &dto.CreatedAt, &dto.UpdatedAt, &dto.LikeCount); err != nil {
			return nil, postgres.MapError(err, "review")
		}

		review, err := dto.ToDomain()
		if err != nil {
			return nil, err
		}

		list = append(list, *review)
	}
	if err := rows.Err(); err != nil {
		return nil, postgres.MapError(err, "review")
	}

	return list, nil
}

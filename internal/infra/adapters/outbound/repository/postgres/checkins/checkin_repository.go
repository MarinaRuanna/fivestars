package checkins

import (
	"context"
	"fmt"
	"time"

	"fivestars/internal/domain"
	"fivestars/internal/infra/adapters/outbound/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type checkinRepository struct {
	pool *pgxpool.Pool
}

func NewCheckinRepository(pool *pgxpool.Pool) domain.CheckinRepository {
	return &checkinRepository{pool: pool}
}

func (r *checkinRepository) Create(ctx context.Context, checkin *domain.Checkin) error {
	if checkin.CreatedAt.IsZero() {
		checkin.CreatedAt = time.Now().UTC()
	}
	_, err := r.pool.Exec(ctx, `
		INSERT INTO checkins (id, user_id, establishment_id, lat, lng, checked_at, created_at)
		VALUES (COALESCE(NULLIF($1, ''), uuid_generate_v4()), $2, $3, $4, $5, $6, $7)
	`, checkin.ID, checkin.UserID, checkin.EstablishmentID, checkin.Lat, checkin.Lng, checkin.CheckedAt, checkin.CreatedAt)
	if err != nil {
		return postgres.MapError(fmt.Errorf("insert checkin: %w", err), "checkin")
	}
	return nil
}

func (r *checkinRepository) ListByUser(ctx context.Context, userID string) ([]domain.Checkin, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, establishment_id, lat, lng, checked_at, created_at
		FROM checkins
		WHERE user_id = $1
		ORDER BY checked_at DESC
	`, userID)
	if err != nil {
		return nil, postgres.MapError(err, "checkin")
	}
	defer rows.Close()

	var list []domain.Checkin
	for rows.Next() {
		var dto CheckinDTO
		err := rows.Scan(&dto.ID, &dto.UserID, &dto.EstablishmentID, &dto.Lat, &dto.Lng, &dto.CheckedAt, &dto.CreatedAt)
		if err != nil {
			return nil, postgres.MapError(err, "checkin")
		}
		list = append(list, *dto.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, postgres.MapError(err, "checkin")
	}
	return list, nil
}

func (r *checkinRepository) FindTodayByUserAndEstablishment(ctx context.Context, userID, establishmentID string) (*domain.Checkin, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, establishment_id, lat, lng, checked_at, created_at
		FROM checkins
		WHERE user_id = $1 AND establishment_id = $2
		AND date_trunc('day', checked_at) = date_trunc('day', now() at time zone 'utc')
		LIMIT 1
	`, userID, establishmentID)
	var dto CheckinDTO
	if err := row.Scan(&dto.ID, &dto.UserID, &dto.EstablishmentID, &dto.Lat, &dto.Lng, &dto.CheckedAt, &dto.CreatedAt); err != nil {
		if postgres.IsNoRows(err) {
			return nil, nil
		}
		return nil, postgres.MapError(err, "checkin")
	}
	return dto.ToDomain(), nil
}

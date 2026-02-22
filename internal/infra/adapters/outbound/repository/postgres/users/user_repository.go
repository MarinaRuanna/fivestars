package users

import (
	"context"
	"errors"

	"fivestars/internal/domain"
	"fivestars/internal/domain/customerror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// userRepository implements domain.UserRepository.
type userRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository returns a new UserRepository.
func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{pool: pool}
}

// Create persiste um novo usuário. O ID e timestamps são gerados pelo banco.
func (r *userRepository) Create(ctx context.Context, u *domain.User) error {
	var userID pgtype.UUID
	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, name, avatar_url, level)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5)
		RETURNING id, created_at, updated_at
	`, u.Email, u.PasswordHash, u.Name, u.AvatarURL, u.Level).Scan(&userID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return customerror.NewConflictError("email já cadastrado")
		}
		return err
	}
	u.ID = uuidToString(userID)
	return nil
}

// GetByID retorna o usuário pelo ID ou nil se não existir.
func (r *userRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	var dto UserDTO
	uuid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}
	err = r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, name, avatar_url, level, created_at, updated_at
		FROM users WHERE id = $1
	`, uuid).Scan(
		&dto.ID, &dto.Email, &dto.PasswordHash, &dto.Name,
		&dto.AvatarURL, &dto.Level, &dto.CreatedAt, &dto.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	user, err := dto.ToDomain()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retorna o usuário pelo email ou nil se não existir.
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dto UserDTO
	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password_hash, name, avatar_url, level, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&dto.ID, &dto.Email, &dto.PasswordHash, &dto.Name,
		&dto.AvatarURL, &dto.Level, &dto.CreatedAt, &dto.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return dto.ToDomain()
}

func parseUUID(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	err := u.Scan(s)
	return u, err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

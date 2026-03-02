package postgres

import (
	"errors"
	"fmt"
	"strings"

	"fivestars/internal/domain/customerror"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type customErrorConstructorFunc func(message string) error

var databaseErrorHandler = map[string]customErrorConstructorFunc{
	"23505": customerror.NewConflictError,   // unique_violation
	"23503": customerror.NewValidationError, // foreign_key_violation
	"22P02": customerror.NewValidationError, // invalid_text_representation
	"23502": customerror.NewValidationError, // not_null_violation
	"22001": customerror.NewValidationError, // string_data_right_truncation
	"22003": customerror.NewValidationError, // numeric_value_out_of_range
	"23514": customerror.NewValidationError, // check_violation
	"57014": customerror.NewValidationError, // query_canceled
}

// IsNoRows reports whether the error means no rows were found.
func IsNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

// MapError translates known PostgreSQL errors to domain errors.
// Unknown database errors are returned wrapped for upstream 500 handling.
func MapError(err error, entity string) error {
	if err == nil {
		return nil
	}

	if IsNoRows(err) {
		return customerror.NewNotFoundError(fmt.Sprintf("%s not found", entity))
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if strings.HasPrefix(pgErr.Code, "08") || pgErr.Code == "57P01" {
			return customerror.NewServiceUnavailableError("database temporarily unavailable")
		}

		if constructor, ok := databaseErrorHandler[pgErr.Code]; ok {
			return constructor(buildErrorMessage(pgErr.Code, entity))
		}
	}

	return fmt.Errorf("database %s error: %w", entity, err)
}

func buildErrorMessage(code, entity string) string {
	switch code {
	case "23505":
		return fmt.Sprintf("%s already exists", entity)
	case "23503":
		return fmt.Sprintf("invalid %s reference", entity)
	case "22P02":
		return fmt.Sprintf("invalid %s value", entity)
	case "23502":
		return fmt.Sprintf("missing required %s field", entity)
	case "22001":
		return fmt.Sprintf("%s value is too long", entity)
	case "22003":
		return fmt.Sprintf("%s value is out of range", entity)
	case "23514":
		return fmt.Sprintf("%s constraint violation", entity)
	case "57014":
		return "database query canceled"
	default:
		return fmt.Sprintf("invalid %s data", entity)
	}
}

package usecases_test

import (
	"errors"
	"fivestars/internal/domain/customerror"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireCustomErrorType(t *testing.T, err error, expected customerror.ErrorType) {
	t.Helper()
	require.Error(t, err)
	var ce *customerror.CustomError
	require.True(t, errors.As(err, &ce))
	require.Equal(t, expected, ce.ErrorType())
}

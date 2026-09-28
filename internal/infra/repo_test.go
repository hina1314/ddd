package infra

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestIsDuplicateKeyError(t *testing.T) {
	require.True(t, IsDuplicateKeyError(&pgconn.PgError{Code: "23505"}))
	require.False(t, IsDuplicateKeyError(&pgconn.PgError{Code: "23503"}))
	require.False(t, IsDuplicateKeyError(errors.New("duplicate text is not a typed database error")))
}

package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	r := NewRandUtil()
	password := r.String(6)
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := r.String(6)

	err = CheckPassword(wrongPassword, hashedPassword)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}

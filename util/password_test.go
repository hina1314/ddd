package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPassword(t *testing.T) {
	password := "correct-horse-battery-staple"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	err = CheckPassword("wrong-password", hashedPassword)

	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}

package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testKey = "01234567890123456789012345678901"

func TestPasetoRoundTrip(t *testing.T) {
	maker, err := NewPasetoMaker(testKey, "issuer", "audience")
	require.NoError(t, err)

	encoded, err := maker.CreateToken(42, time.Minute)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(encoded)
	require.NoError(t, err)
	require.Equal(t, int64(42), payload.UserID)
	require.Equal(t, "issuer", payload.Issuer)
	require.Equal(t, "audience", payload.Audience)
}

func TestPasetoRejectsWrongAudience(t *testing.T) {
	issuer, err := NewPasetoMaker(testKey, "issuer", "audience-a")
	require.NoError(t, err)
	verifier, err := NewPasetoMaker(testKey, "issuer", "audience-b")
	require.NoError(t, err)

	encoded, err := issuer.CreateToken(42, time.Minute)
	require.NoError(t, err)

	_, err = verifier.VerifyToken(encoded)
	require.ErrorIs(t, err, ErrInvalidToken)
}

func TestPayloadRejectsExpiredAndInvalidTimes(t *testing.T) {
	now := time.Now()
	payload := &Payload{ID: "id", UserID: 1, Issuer: "issuer", Audience: "audience", IssuedAt: now.Add(-time.Hour), ExpiredAt: now.Add(-time.Minute)}
	require.ErrorIs(t, payload.Valid("issuer", "audience"), ErrExpiredToken)

	payload.ExpiredAt = payload.IssuedAt.Add(-time.Second)
	require.ErrorIs(t, payload.Valid("issuer", "audience"), ErrInvalidToken)
}

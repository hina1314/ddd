package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// different types of error returned by the VerifyToken function
var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains payload data of the token
type Payload struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	Issuer    string    `json:"issuer"`
	Audience  string    `json:"audience"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates an access-token payload without embedding user PII.
func NewPayload(userID int64, duration time.Duration, issuer, audience string) (*Payload, error) {
	if userID <= 0 || duration <= 0 || issuer == "" || audience == "" {
		return nil, ErrInvalidToken
	}
	now := time.Now()
	payload := &Payload{
		ID:        uuid.NewString(),
		UserID:    userID,
		Issuer:    issuer,
		Audience:  audience,
		IssuedAt:  now,
		ExpiredAt: now.Add(duration),
	}
	return payload, nil
}

// Valid checks the token identity, intended recipients and timestamps.
func (payload *Payload) Valid(issuer, audience string) error {
	now := time.Now()
	if payload.ID == "" || payload.UserID <= 0 || payload.Issuer != issuer || payload.Audience != audience {
		return ErrInvalidToken
	}
	if payload.IssuedAt.IsZero() || payload.ExpiredAt.IsZero() ||
		payload.IssuedAt.After(now.Add(time.Minute)) || !payload.ExpiredAt.After(payload.IssuedAt) {
		return ErrInvalidToken
	}
	if !payload.ExpiredAt.After(now) {
		return ErrExpiredToken
	}
	return nil
}

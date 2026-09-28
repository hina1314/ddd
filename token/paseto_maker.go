package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto/v2"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
	issuer       string
	audience     string
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey, issuer, audience string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d bytes", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
		issuer:       issuer,
		audience:     audience,
	}
	return maker, nil
}

// CreateToken creates a new token for a specific userId and duration
func (maker *PasetoMaker) CreateToken(userID int64, duration time.Duration) (string, error) {
	payload, err := NewPayload(userID, duration, maker.issuer, maker.audience)
	if err != nil {
		return "", err
	}

	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid(maker.issuer, maker.audience)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

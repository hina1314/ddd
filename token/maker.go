package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// CreateToken creates a new access token for a user.
	CreateToken(userID int64, duration time.Duration) (string, error)
	//VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}

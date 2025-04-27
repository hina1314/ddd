package entity

import (
	"time"
)

type Hotel struct {
	ID             int64
	MerchantID     int64
	Name           string
	Address        string
	City           string
	Country        string
	Phone          string
	Email          string
	PoliceCode     string
	PoliceAuthCode string
	PoliceSign     string
	CreatedAt      time.Time
}

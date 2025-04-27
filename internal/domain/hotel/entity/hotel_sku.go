package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type HotelSku struct {
	ID              int64
	HotelID         int64
	RoomID          int64
	Name            string
	SalesPrice      decimal.Decimal
	TicketPrice     decimal.Decimal
	TicketStatus    string
	RefundStatus    bool
	RefundAudit     string
	RefundCondition string // JSON
	Status          string
	hotel           Hotel
	CreatedAt       time.Time
}

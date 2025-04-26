package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type HotelSku struct {
	ID              uint
	HotelID         uint
	RoomID          uint
	Name            string
	SalesPrice      decimal.Decimal
	TicketPrice     decimal.Decimal
	TicketStatus    string
	RefundStatus    string
	RefundAudit     string
	RefundCondition string // JSON
	Status          string
	CreatedAt       time.Time
}

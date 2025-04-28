package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type HotelSku struct {
	ID              int64
	MerchantID      int64
	HotelID         int64
	RoomTypeID      int64
	Name            string
	SalesPrice      decimal.Decimal
	TicketPrice     decimal.Decimal
	TicketStatus    string
	RefundStatus    bool
	RefundAudit     string
	RefundCondition string // JSON
	Status          string
	Hotel           Hotel
	CreatedAt       time.Time
}

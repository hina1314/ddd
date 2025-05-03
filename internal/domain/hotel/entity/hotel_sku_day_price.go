package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type HotelSkuDayPrice struct {
	ID           int64
	HotelID      int64
	SkuID        int64
	RoomTypeID   int64
	Date         time.Time
	MarketPrice  decimal.Decimal
	SalePrice    decimal.Decimal
	TicketPrice  decimal.Decimal
	TicketStatus bool
	CreatedAt    time.Time
}

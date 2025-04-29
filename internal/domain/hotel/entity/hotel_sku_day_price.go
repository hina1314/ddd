package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type SkuDayPrice struct {
	HotelId      int64
	SkuId        int64
	RoomTypeId   int64
	Date         string
	MarketPrice  decimal.Decimal
	SalePrice    decimal.Decimal
	TicketPrice  decimal.Decimal
	TicketStatus bool
	CreatedAt    time.Time
}

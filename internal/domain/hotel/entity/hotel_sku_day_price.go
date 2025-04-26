package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type DatePrice struct {
	HotelId      int64
	SkuId        int64
	RoomTypeId   int64
	Date         string
	MarketPrice  decimal.Decimal
	SalePrice    decimal.Decimal
	TicketPrice  decimal.Decimal
	TicketStatus int8
	CreatedAt    time.Time
}

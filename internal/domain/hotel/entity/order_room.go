package entity

import (
	"github.com/shopspring/decimal"
	"time"
)

type OrderRoomStatus string

const (
	OrderRoomStatusInit OrderRoomStatus = "init"
)

type OrderRoom struct {
	ID         uint
	HotelID    int64
	MerchantID int64
	RoomTypeID int64
	RoomItemID int64
	Price      decimal.Decimal
	Status     OrderRoomStatus
	ExpireTime time.Time
	CreatedAt  time.Time
}

func NewOrderRoom(userId int64, hotelSku HotelSku, totalPrice decimal.Decimal, totalNumber, totalPayTicket int) *OrderRoom {
	return &OrderRoom{
		HotelID:    hotelSku.HotelID,
		MerchantID: hotelSku.hotel.MerchantID,
		Status:     OrderRoomStatusInit,
		ExpireTime: time.Now().Add(1800 * time.Second),
		CreatedAt:  time.Time{},
	}
}

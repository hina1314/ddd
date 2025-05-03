package entity

import (
	"github.com/shopspring/decimal"
	"study/internal/domain/hotel/entity"
	"time"
)

type OrderRoomStatus int16

const (
	OrderRoomStatusInit OrderRoomStatus = 0
)

type OrderRoom struct {
	ID         int64
	OrderID    int64
	HotelID    int64
	MerchantID int64
	RoomTypeID int64
	RoomItemID int64
	Price      decimal.Decimal
	Status     OrderRoomStatus
	CreatedAt  time.Time
}

func NewOrderRoom(userId int64, hotelSku entity.HotelSku, totalPrice decimal.Decimal, totalNumber, totalPayTicket int) *OrderRoom {
	return &OrderRoom{
		HotelID:    hotelSku.HotelID,
		MerchantID: hotelSku.Hotel.MerchantID,
		Status:     OrderRoomStatusInit,
		CreatedAt:  time.Time{},
	}
}

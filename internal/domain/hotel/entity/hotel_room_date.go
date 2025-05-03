package entity

import (
	"time"
)

type RoomDateStatus int16

const (
	RoomDateStatusOrderLock   RoomDateStatus = 0
	RoomDateStatusWaitCheckin RoomDateStatus = 1
	RoomDateStatusCheckin     RoomDateStatus = 2
	RoomDateStatusCheckout    RoomDateStatus = 3
	RoomDateStatusLock        RoomDateStatus = 4
)

type HotelRoomDate struct {
	ID         int64
	OrderID    int64
	HotelID    int64
	RoomItemID int64
	Date       time.Time
	Status     RoomDateStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (rd *HotelRoomDate) MarkAsCheckin() {
	rd.Status = RoomDateStatusCheckin
}

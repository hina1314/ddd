package entity

import (
	"time"
)

type RoomDateStatus string

const (
	RoomDateStatusAvailable RoomDateStatus = "available"
	RoomDateStatusLocked    RoomDateStatus = "locked"
	RoomDateStatusCheckin   RoomDateStatus = "checkin"
)

type RoomDate struct {
	ID             uint
	OrderID        uint
	HotelID        uint
	RoomID         uint
	RoomInstanceID uint
	Date           time.Time
	Status         RoomDateStatus
	CreatedAt      time.Time
}

func (rd *RoomDate) MarkAsCheckin() {
	rd.Status = RoomDateStatusCheckin
}

package co

import "time"

type contact struct {
	Name  string
	Phone string
}

type RoomContact struct {
	Guests []contact
}

type CreateOrderCommand struct {
	SkuID     string
	StartDate time.Time
	EndDate   time.Time
	Number    int
	PriceType int8
	PayType   string
	Contact   []RoomContact
}

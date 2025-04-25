package entity

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusInit     OrderStatus = "init"
	OrderStatusPaid     OrderStatus = "paid"
	OrderStatusCheckin  OrderStatus = "checkin"
	OrderStatusCheckout OrderStatus = "checkout"
	OrderStatusRefunded OrderStatus = "refunded"
)

type Order struct {
	ID                   uint
	Ordersn              string
	CustomerID           uint
	HotelID              uint
	MerchantID           uint
	ProductCategory      string
	TotalPrice           float64
	TotalNumber          int
	TotalPayTicket       int
	TotalRefundedAmount  float64
	TotalRefundedTickets int
	TotalRefundedNumber  int
	AllowRefund          string
	Status               OrderStatus
	ExpireTime           *time.Time
	CreatedAt            time.Time
}

func (o *Order) CanRefund() bool {
	return o.AllowRefund == "yes"
}

func (o *Order) UpdateRefundStats(amount float64, tickets, number int) {
	o.TotalRefundedAmount += amount
	o.TotalRefundedTickets += tickets
	o.TotalRefundedNumber += number
}

func (o *Order) MarkAsRefunded() {
	o.Status = OrderStatusRefunded
}

func (o *Order) MarkAsCheckin() {
	o.Status = OrderStatusCheckin
}

func (o *Order) MarkAsCheckout() {
	o.Status = OrderStatusCheckout
}

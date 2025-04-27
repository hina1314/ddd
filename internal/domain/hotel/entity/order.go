package entity

import (
	"fmt"
	"github.com/shopspring/decimal"
	"study/util"
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
	OrderSn              string
	UserId               int64
	HotelID              int64
	MerchantID           int64
	ProductCategory      string
	TotalPrice           decimal.Decimal
	TotalNumber          int
	TotalPayTicket       int
	TotalRefundedAmount  decimal.Decimal
	TotalRefundedTickets int
	TotalRefundedNumber  int
	AllowRefund          bool
	Status               OrderStatus
	ExpireTime           *time.Time
	CreatedAt            time.Time
}

func NewOrder(userId, hotelId, merchantId int64, totalPrice decimal.Decimal, totalNumber, totalPayTicket int) *Order {
	return &Order{
		OrderSn:              orderSn("LJ"),
		UserId:               userId,
		HotelID:              hotelId,
		MerchantID:           merchantId,
		ProductCategory:      "",
		TotalPrice:           totalPrice,
		TotalNumber:          totalNumber,
		TotalPayTicket:       totalPayTicket,
		TotalRefundedAmount:  decimal.Zero,
		TotalRefundedTickets: 0,
		TotalRefundedNumber:  0,
		AllowRefund:          true,
		Status:               OrderStatusInit,
		ExpireTime:           nil,
		CreatedAt:            time.Time{},
	}
}

func (o *Order) CanRefund() bool {
	return o.AllowRefund
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

func orderSn(pre string) string {
	now := time.Now().Format("20060102150405") // Go的时间模板
	r := util.NewRandUtil()
	num := r.Int(1000, 9999)
	return fmt.Sprintf("%v%v%v", pre, now, num)
}

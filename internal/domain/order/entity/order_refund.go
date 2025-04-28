package entity

import (
	"time"
)

type RefundStatus string

const (
	RefundStatusPending  RefundStatus = "pending"
	RefundStatusApproved RefundStatus = "approved"
	RefundStatusRejected RefundStatus = "rejected"
)

type OrderRefund struct {
	ID            uint
	OrderID       uint
	OrderRoomID   uint
	Amount        float64
	TicketCount   int
	Reason        string
	ReasonType    string
	Status        RefundStatus
	RefundSn      string
	RefundTime    *time.Time
	ThirdRefundID string
	CreatedAt     time.Time
}

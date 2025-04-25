package entity

import (
	"time"
)

type UserPlanStatus string

const (
	UserPlanStatusInit       UserPlanStatus = "init"
	UserPlanStatusCheckWait  UserPlanStatus = "check_wait"
	UserPlanStatusCheckin    UserPlanStatus = "checkin"
	UserPlanStatusCheckout   UserPlanStatus = "checkout"
	UserPlanStatusRefundWait UserPlanStatus = "refund_wait"
	UserPlanStatusRefunded   UserPlanStatus = "refunded"
)

type UserPlan struct {
	ID             uint
	OrderID        uint
	OrderRoomID    uint
	RoomInstanceID uint
	CustomerID     uint
	Name           string
	Phone          string
	Title          string
	StartDate      time.Time
	EndDate        time.Time
	Status         UserPlanStatus
	CreatedAt      time.Time
}

func (p *UserPlan) MarkAsCheckin() {
	p.Status = UserPlanStatusCheckin
}

func (p *UserPlan) MarkAsCheckout() {
	p.Status = UserPlanStatusCheckout
}

func (p *UserPlan) MarkAsRefundWait() {
	p.Status = UserPlanStatusRefundWait
}

func (p *UserPlan) MarkAsRefunded() {
	p.Status = UserPlanStatusRefunded
}

package entity

import "time"

type UserPlanStatus int16

const (
	UserPlanStatusInit UserPlanStatus = 0
)

type UserPlan struct {
	ID         int64
	OrderID    int64
	RoomItemID int64
	Phone      string
	Name       string
	Start      time.Time
	End        time.Time
	Status     UserPlanStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

type Contact struct {
	Name  string
	Phone string
}

type RoomContact struct {
	RoomItemID int64     // 指定的房间标识
	Guests     []Contact // 这个房间下的所有联系人
}

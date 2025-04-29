package service

import (
	"errors"
	"study/internal/assemble"
	"study/internal/domain/user"
	"time"
)

type UserPlanService struct{}

func NewUserPlanService() *UserPlanService {
	return &UserPlanService{}
}

// BuildPlans 构建用户行程数据
func (s *UserPlanService) BuildPlans(orderID int64, roomItemIDs []int64, contacts []assemble.RoomContact, start, end time.Time) ([]user.UserPlan, error) {
	if len(roomItemIDs) != len(contacts) {
		return nil, errors.New("room and contact count mismatch")
	}

	var plans []user.UserPlan

	for i, roomContact := range contacts {
		for _, c := range roomContact.Guests {
			plans = append(plans, user.UserPlan{
				OrderId:    orderID,
				RoomItemID: roomItemIDs[i],
				Phone:      c.Phone,
				Name:       c.Name,
				Start:      start,
				End:        end,
				Status:     user.UserPlanStatusInit,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		}
	}

	return plans, nil
}

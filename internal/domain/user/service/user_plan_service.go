package service

import (
	"context"
	"strings"
	"study/internal/domain/user/entity"
	"study/internal/domain/user/repository"
	"study/util/errors"
	"time"
)

type UserPlanService struct {
	userPlanRepo repository.UserPlanRepository
}

func NewUserPlanService(userPlan repository.UserPlanRepository) *UserPlanService {
	return &UserPlanService{
		userPlanRepo: userPlan,
	}
}

func (s *UserPlanService) CreateUserPlans(ctx context.Context, orderID int64, roomItemIDs []int64, roomContacts []entity.RoomContact, start, end time.Time) error {
	if len(roomContacts) != len(roomItemIDs) {
		return errors.New("xxx", "invalid_contact_count room number mismatch")
	}
	plans, err := s.BuildPlans(orderID, roomItemIDs, roomContacts, start, end)
	if err != nil {
		return err
	}
	return s.userPlanRepo.AddUserPlan(ctx, plans)
}

// BuildPlans 构建用户行程数据
func (s *UserPlanService) BuildPlans(orderID int64, roomItemIDs []int64, roomContacts []entity.RoomContact, start, end time.Time) ([]entity.UserPlan, error) {
	var plans []entity.UserPlan

	for i, roomContact := range roomContacts {
		for _, c := range roomContact.Guests {
			plans = append(plans, entity.UserPlan{
				OrderID:    orderID,
				RoomItemID: roomItemIDs[i],
				Phone:      c.Phone,
				Name:       c.Name,
				Start:      start,
				End:        end,
				Status:     entity.UserPlanStatusInit,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		}
	}

	return plans, nil
}

func (s *UserPlanService) CheckBookingConflicts(ctx context.Context, start, end time.Time, roomContacts []entity.RoomContact) error {
	for _, roomContact := range roomContacts {
		for _, contact := range roomContact.Guests {
			contact.Phone = removePhoneSpaces(contact.Phone)
			if err := s.userPlanRepo.CheckBookingConflicts(ctx, start, end, contact.Phone); err != nil {
				return err
			}
		}
	}
	return nil
}

func removePhoneSpaces(phone string) string {
	return strings.ReplaceAll(phone, " ", "")
}

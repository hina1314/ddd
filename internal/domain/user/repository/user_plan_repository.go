package repository

import (
	"context"
	"study/internal/domain/user/entity"
	"time"
)

type UserPlanRepository interface {
	CheckBookingConflicts(ctx context.Context, start, end time.Time, phone string) error
	AddUserPlan(ctx context.Context, userPlan []entity.UserPlan) error
}

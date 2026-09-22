package user

import (
	"context"
	"study/db/model"
	"study/internal/domain/user/entity"
	"study/internal/domain/user/repository"
	"time"
)

type UserPlanRepositoryImpl struct {
	db model.TxManager
}

func NewUserPlanRepo(db model.TxManager) repository.UserPlanRepository {
	return &UserPlanRepositoryImpl{
		db: db,
	}
}

func (u *UserPlanRepositoryImpl) CheckBookingConflicts(ctx context.Context, start, end time.Time, phone string) error {
	//q := u.db.Querier(ctx)

	//count, err := q.CheckBookingConflicts(ctx, model.CheckBookingConflictsParams{
	//	Phone:     phone,
	//	StartDate: start,
	//	EndDate:   end,
	//})
	//if err != nil {
	//	return errors.New("check_booking_conflicts_error", err.Error())
	//}
	//
	//if count > 0 {
	//	return errors.New(errors.ErrBookingConflict, "user has conflicting bookings for the specified dates")
	//}

	return nil
}

func (u *UserPlanRepositoryImpl) AddUserPlan(ctx context.Context, userPlans []entity.UserPlan) error {
	// q := u.db.Querier(ctx)

	// Prepare arrays for bulk insert
	orderIDs := make([]int64, len(userPlans))
	roomItemIDs := make([]int64, len(userPlans))
	phones := make([]string, len(userPlans))
	names := make([]string, len(userPlans))
	startDates := make([]time.Time, len(userPlans))
	endDates := make([]time.Time, len(userPlans))
	statuses := make([]int16, len(userPlans))
	createdAts := make([]time.Time, len(userPlans))
	updatedAts := make([]time.Time, len(userPlans))

	for i, plan := range userPlans {
		orderIDs[i] = plan.OrderID
		roomItemIDs[i] = plan.RoomItemID
		phones[i] = plan.Phone
		names[i] = plan.Name
		startDates[i] = plan.Start
		endDates[i] = plan.End
		statuses[i] = int16(plan.Status)
		createdAts[i] = plan.CreatedAt
		updatedAts[i] = plan.UpdatedAt
	}

	//err := q.AddUserPlan(ctx, model.AddUserPlanParams{
	//	OrderID:    orderIDs,
	//	RoomItemID: roomItemIDs,
	//	Phone:      phones,
	//	Name:       names,
	//	StartDate:  startDates,
	//	EndDate:    endDates,
	//	Status:     statuses,
	//	CreatedAt:  createdAts,
	//	UpdatedAt:  updatedAts,
	//})
	//if err != nil {
	//	return errors.New("add_user_plan_error", "failed to add user plans: "+err.Error())
	//}

	return nil
}

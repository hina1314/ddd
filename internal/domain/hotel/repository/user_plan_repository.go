package repository

type UserPlanRepository interface {
	CountConflictingPlans(phone, start, end string) (int, error)
}

package user

import (
	"context"
	"time"
)

// UserRepository 用户实体仓储接口
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
	Save(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*User, int, error)
	CheckBookingConflicts(ctx context.Context, start, end time.Time, phone string) error
	AddUserPlan(ctx context.Context, userPlan []UserPlan) error
}

// UserAccountRepository 用户账户仓储接口
type UserAccountRepository interface {
	GetByID(ctx context.Context, id int64) (*UserAccount, error)
	GetByUserID(ctx context.Context, userID int64) (*UserAccount, error)
	Save(ctx context.Context, account *UserAccount) error
	Update(ctx context.Context, account *UserAccount) error
	Delete(ctx context.Context, id int64) error
}

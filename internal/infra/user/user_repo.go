package user

import (
	"context"
	"database/sql"
	"study/db/model"
	"study/internal/domain/user/entity"
	"study/internal/domain/user/repository"
	"study/internal/infra"
	"study/util/errors"
	"time"
)

type UserRepositoryImpl struct {
	db model.TxManager
}

func NewUserRepository(db model.TxManager) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// **通用转换方法：model.User → user.User**
func (r *UserRepositoryImpl) toDomain(u model.User) (*entity.User, error) {
	emailVO := entity.NewEmail(u.Email.String)
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}

	return &entity.User{
		ID:        u.ID,
		Phone:     u.Phone,
		Username:  u.Username,
		Email:     emailVO,
		Password:  u.Password,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	u, err := r.db.Querier(ctx).GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	u, err := r.db.Querier(ctx).GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	u, err := r.db.Querier(ctx).GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	u, err := r.db.Querier(ctx).GetUserByEmail(ctx, toNullString(email))
	if err != nil {
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) Save(ctx context.Context, u *entity.User) error {
	q := r.db.Querier(ctx)
	arg := model.CreateUserParams{
		Phone:    u.Phone,
		Email:    emailToNullString(u.Email),
		Username: u.Username,
		Password: u.Password,
	}

	result, err := q.CreateUser(ctx, arg)
	if err != nil {
		if infra.IsDuplicateKeyError(err) {
			return errors.New(errors.ErrUserAlreadyExists, "User already exists")
		}
		return err
	}
	u.ID = result.ID
	u.Account.UserID = result.ID
	u.CreatedAt = result.CreatedAt
	u.UpdatedAt = result.UpdatedAt

	err = q.CreateUserAccount(ctx, model.CreateUserAccountParams{
		UserID:        u.ID,
		FrozenBalance: u.Account.FrozenBalance.Amount,
		Balance:       u.Account.Balance.Amount,
	})
	if err != nil {
		if infra.IsDuplicateKeyError(err) {
			return errors.New(errors.ErrUserAlreadyExists, "UserAccount already exists")
		}
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, u *entity.User) error {
	arg := model.UpdateUserParams{
		ID:       u.ID,
		Phone:    u.Phone,
		Email:    emailToNullString(u.Email),
		Username: u.Username,
		Password: u.Password,
	}

	return r.db.Querier(ctx).UpdateUser(ctx, arg)
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.Querier(ctx).DeleteUser(ctx, id)
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.User, int, error) {
	users, err := r.db.Querier(ctx).ListUsers(ctx, model.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.User, 0, len(users))
	for _, u := range users {
		domainUser, err := r.toDomain(u)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, domainUser)
	}

	// **修复 CountUsers 需要查询数据库**
	count, err := r.db.Querier(ctx).CountUsers(ctx)
	if err != nil {
		return nil, 0, err
	}
	return result, int(count), nil
}

func emailToNullString(email entity.Email) sql.NullString {
	if email.IsNil() {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: email.String(), Valid: true}
}

func toNullString(str string) sql.NullString {
	if str == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: str, Valid: true}
}

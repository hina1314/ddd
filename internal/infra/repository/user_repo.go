package repository

import (
	"context"
	"database/sql"
	strErr "errors"
	"study/db/model"
	"study/internal/domain/user"
	"study/util/errors"
	"time"
)

type UserRepositoryImpl struct {
	db model.TxManager // 使用 sqlc 生成的 Queries
}

func NewUserRepository(store model.TxManager) user.UserRepository {
	return &UserRepositoryImpl{
		db: store,
	}
}

func (r *UserRepositoryImpl) getQuerier(ctx context.Context) model.Querier {
	if tx, ok := ctx.Value(model.TxKey{}).(model.Tx); ok {
		return tx
	}
	return r.db
}

// **通用转换方法：model.User → user.User**
func (r *UserRepositoryImpl) toDomain(u model.User) (*user.User, error) {
	emailVO := user.NewEmail(u.Email)

	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}

	return &user.User{
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

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	u, err := r.getQuerier(ctx).GetUserByUsername(ctx, username)
	if err != nil {
		if strErr.Is(err, sql.ErrNoRows) {
			return nil, errors.New(errors.ErrUserNotFound, "User not found")
		}
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	u, err := r.getQuerier(ctx).GetUserByPhone(ctx, phone)
	if err != nil {
		if strErr.Is(err, sql.ErrNoRows) {
			return nil, errors.New(errors.ErrUserNotFound, "User not found")
		}
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*user.User, error) {
	u, err := r.getQuerier(ctx).GetUserByID(ctx, id)
	if err != nil {
		if strErr.Is(err, sql.ErrNoRows) {
			return nil, errors.New(errors.ErrUserNotFound, "User not found")
		}
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	u, err := r.getQuerier(ctx).GetUserByEmail(ctx, email)
	if err != nil {
		if strErr.Is(err, sql.ErrNoRows) {
			return nil, errors.New(errors.ErrUserNotFound, "User not found")
		}
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) Save(ctx context.Context, u *user.User) error {
	arg := model.CreateUserParams{
		Phone:    u.Phone,
		Email:    u.Email.String(),
		Username: u.Username,
		Password: u.Password,
	}

	result, err := r.getQuerier(ctx).CreateUser(ctx, arg)
	if err != nil {
		if isDuplicateKeyError(err) {
			return errors.New(errors.ErrUserAlreadyExists, "User already exists")
		}
		return err
	}
	u.ID = result.ID
	u.CreatedAt = result.CreatedAt
	u.UpdatedAt = result.UpdatedAt
	return nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, u *user.User) error {
	arg := model.UpdateUserParams{
		ID:       u.ID,
		Phone:    u.Phone,
		Email:    u.Email.String(),
		Username: u.Username,
		Password: u.Password,
	}
	return r.getQuerier(ctx).UpdateUser(ctx, arg)
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.getQuerier(ctx).DeleteUser(ctx, id)
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*user.User, int, error) {
	users, err := r.getQuerier(ctx).ListUsers(ctx, model.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*user.User, 0, len(users))
	for _, u := range users {
		domainUser, err := r.toDomain(u)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, domainUser)
	}

	// **修复 CountUsers 需要查询数据库**
	count, err := r.getQuerier(ctx).CountUsers(ctx)
	if err != nil {
		return nil, 0, err
	}
	return result, int(count), nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"study/db/model"
	"study/internal/domain/user"
	"time"
)

type UserRepositoryImpl struct {
	db model.Store // 使用 sqlc 生成的 Queries
}

func NewUserRepository(store model.Store) user.UserRepository {
	return &UserRepositoryImpl{
		db: store,
	}
}

// **通用转换方法：model.User → user.User**
func (r *UserRepositoryImpl) toDomain(u model.User) (*user.User, error) {
	emailVO, err := user.NewEmail(u.Email)
	if err != nil {
		return nil, err
	}

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
	u, err := r.db.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 用户不存在
		}
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*user.User, error) {
	u, err := r.db.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	u, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
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
	result, err := r.db.CreateUser(ctx, arg)
	if err != nil {
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
	return r.db.UpdateUser(ctx, arg)
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.db.DeleteUser(ctx, id)
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*user.User, int, error) {
	users, err := r.db.ListUsers(ctx, model.ListUsersParams{
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
	count, err := r.db.CountUsers(ctx)
	if err != nil {
		return nil, 0, err
	}
	return result, int(count), nil
}

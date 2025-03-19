package repository

import (
	"context"
	"database/sql"
	"study/db/model"
	"study/internal/domain/user"
)

type UserRepositoryImpl struct {
	db model.Store // 使用 sqlc 生成的 Queries
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &UserRepositoryImpl{
		db: model.NewStore(db), // 初始化 sqlc 的 Queries
	}
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*user.User, error) {
	//u, err := r.db.GetUserByID(ctx, int32(id)) // 假设 sqlc 生成的方法
	//if err != nil {
	//	if err == sql.ErrNoRows {
	//		return nil, nil // 未找到用户返回 nil
	//	}
	//	return nil, err
	//}
	//return &user.User{
	//	ID:        int(u.ID),
	//	Phone:     u.Phone,
	//	Username:  u.Username,
	//	Email:     u.Email.String(),
	//	Password:  u.Password,
	//	CreatedAt: u.CreatedAt,
	//}, nil
	panic("implement me")
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	//u, err := r.db.GetUserByUsername(ctx, username) // 假设 sqlc 生成的方法
	//if err != nil {
	//	if err == sql.ErrNoRows {
	//		return nil, nil
	//	}
	//	return nil, err
	//}
	//return &user.User{
	//	ID:        int(u.ID),
	//	Phone:     u.Phone,
	//	Username:  u.Username,
	//	Email:     u.Email.String(),
	//	Password:  u.Password,
	//	CreatedAt: u.CreatedAt,
	//}, nil
	panic("implement me")
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	//u, err := r.db.GetUserByEmail(ctx, email) // 假设 sqlc 生成的方法
	//if err != nil {
	//	if err == sql.ErrNoRows {
	//		return nil, nil
	//	}
	//	return nil, err
	//}
	//return &user.User{
	//	ID:        int(u.ID),
	//	Phone:     u.Phone,
	//	Username:  u.Username,
	//	Email:     u.Email.String(),
	//	Password:  u.Password,
	//	CreatedAt: u.CreatedAt,
	//}, nil
	panic("implement me")
}

func (r *UserRepositoryImpl) Save(ctx context.Context, user *user.User) error {
	arg := model.CreateUserParams{
		Phone:    user.Phone,
		Email:    user.Email.String(),
		Username: user.Username,
		Password: user.Password,
	}
	u, err := r.db.CreateUser(ctx, arg)
	if err != nil {
		return err
	}
	user.ID = u.ID // 更新 user 的 ID
	user.CreatedAt = u.CreatedAt
	return nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *user.User) error {
	//arg := model.UpdateUserParams{
	//	ID:       int32(user.ID),
	//	Phone:    user.Phone,
	//	Email:    sql.NullString{String: user.Email, Valid: user.Email != ""},
	//	Username: user.Username,
	//	Password: user.Password,
	//}
	//return r.db.UpdateUser(ctx, arg) // 假设 sqlc 生成的方法
	panic("implement me")
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	//return r.db.DeleteUser(ctx, id) // 假设 sqlc 生成的方法
	panic("implement me")
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*user.User, int, error) {
	//users, err := r.db.ListUsers(ctx, model.ListUsersParams{
	//	Limit:  int32(limit),
	//	Offset: int32(offset),
	//})
	//if err != nil {
	//	return nil, 0, err
	//}
	//
	//result := make([]*user.User, len(users))
	//for i, u := range users {
	//	result[i] = &user.User{
	//		ID:        int(u.ID),
	//		Phone:     u.Phone,
	//		Username:  u.Username,
	//		Email:     u.Email.String(),
	//		Password:  u.Password,
	//		CreatedAt: u.CreatedAt,
	//	}
	//}
	//
	//count, err := r.db.CountUsers(ctx) // 假设 sqlc 生成的方法
	//if err != nil {
	//	return nil, 0, err
	//}
	//return result, int(count), nil
	panic("implement me")
}

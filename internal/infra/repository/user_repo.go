package repository

import (
	"context"
	"database/sql"
	"study/internal/domain/user"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (u UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepositoryImpl) Save(ctx context.Context, user *user.User) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepositoryImpl) Update(ctx context.Context, user *user.User) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*user.User, int, error) {
	//TODO implement me
	panic("implement me")
}

package repository

import (
	"context"
	"database/sql"
	"study/internal/domain/user"
)

type UserAccountRepositoryImpl struct {
	db *sql.DB
}

func NewUserAccountRepository(db *sql.DB) user.UserAccountRepository {
	return &UserAccountRepositoryImpl{}
}

func (u UserAccountRepositoryImpl) GetByID(ctx context.Context, id int64) (*user.UserAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserAccountRepositoryImpl) GetByUserID(ctx context.Context, userID string) (*user.UserAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserAccountRepositoryImpl) Save(ctx context.Context, account *user.UserAccount) error {
	//TODO implement me
	panic("implement me")
}

func (u UserAccountRepositoryImpl) Update(ctx context.Context, account *user.UserAccount) error {
	//TODO implement me
	panic("implement me")
}

func (u UserAccountRepositoryImpl) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

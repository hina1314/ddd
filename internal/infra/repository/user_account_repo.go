package repository

import (
	"context"
	"database/sql"
	"study/db/model"
	"study/internal/domain/user"
)

type UserAccountRepositoryImpl struct {
	db model.Store
}

func NewUserAccountRepository(db *sql.DB) user.UserAccountRepository {
	return &UserAccountRepositoryImpl{
		db: model.NewStore(db),
	}
}

func (u UserAccountRepositoryImpl) GetByID(ctx context.Context, id int64) (*user.UserAccount, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserAccountRepositoryImpl) GetByUserID(ctx context.Context, userID int64) (*user.UserAccount, error) {
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

func (u UserAccountRepositoryImpl) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

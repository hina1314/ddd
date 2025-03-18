package repository

import (
	"context"
	"study/internal/domain/user"
	"study/internal/svc"
)

type UserAccountRepositoryImpl struct {
	ctxSvc *svc.ServiceContext
}

func NewUserAccountRepository(ctxSvc *svc.ServiceContext) user.UserAccountRepository {
	return &UserAccountRepositoryImpl{ctxSvc: ctxSvc}
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

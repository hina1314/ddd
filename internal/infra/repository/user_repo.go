package repository

import (
	"context"
	"study/db/model"
	"study/internal/domain/user"
	"study/internal/svc"
)

type UserRepositoryImpl struct {
	ctxSvc *svc.ServiceContext
}

func NewUserRepository(ctxSvc *svc.ServiceContext) user.UserRepository {
	return &UserRepositoryImpl{ctxSvc: ctxSvc}
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
	arg := model.CreateUserParams{
		Phone:    user.Phone,
		Email:    user.Email,
		Username: user.Username,
		Password: user,
	}
	u.ctxSvc.Db.CreateUser(ctx, arg)
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

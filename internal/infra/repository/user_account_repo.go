package repository

import (
	"context"
	"study/db/model"
	"study/internal/domain/user"
)

type UserAccountRepositoryImpl struct {
	db model.TxManager
}

func NewUserAccountRepository(store model.TxManager) user.UserAccountRepository {
	return &UserAccountRepositoryImpl{
		db: store,
	}
}

func (r *UserAccountRepositoryImpl) getQuerier(ctx context.Context) model.Querier {
	if tx, ok := ctx.Value(model.TxKey{}).(model.Tx); ok {
		return tx
	}
	return r.db
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

	arg := model.CreateAccountParams{
		UserID:        account.UserID,
		FrozenBalance: account.FrozenBalance.String(),
		Balance:       account.Balance.String(),
	}
	newAcc, err := u.getQuerier(ctx).CreateAccount(ctx, arg)
	if err != nil {
		return err
	}

	account.ID = newAcc.ID
	account.CreatedAt = newAcc.CreatedAt
	account.UpdatedAt = newAcc.UpdatedAt
	return nil
}

func (u UserAccountRepositoryImpl) Update(ctx context.Context, account *user.UserAccount) error {
	//TODO implement me
	panic("implement me")
}

func (u UserAccountRepositoryImpl) Delete(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

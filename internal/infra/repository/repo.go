package repository

import (
	"context"
	"study/db/model"
)

func (r *UserRepositoryImpl) getQuerier(ctx context.Context) model.Querier {
	if tx, ok := ctx.Value(model.TxKey{}).(model.Tx); ok {
		return tx
	}
	return r.db
}

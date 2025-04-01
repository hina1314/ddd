package repository

import (
	"context"
	"study/db/model"
)

type txKey struct{}

func (r *UserRepositoryImpl) getQuerier(ctx context.Context) model.Querier {
	if tx, ok := ctx.Value(txKey{}).(model.Tx); ok {
		return tx
	}
	return r.db
}

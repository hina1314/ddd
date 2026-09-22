package model

import (
	"context"
	"database/sql"
)

type txKey struct{} // 唯一的事务键

// TxManager 定义事务管理接口
type TxManager interface {
	Querier(ctx context.Context) Querier
	Begin(ctx context.Context) (Tx, context.Context, error)
}

// Tx 定义事务操作接口
type Tx interface {
	Querier // 继承sqlc生成的查询接口
	Commit() error
	Rollback() error
}

// SQLStore provides all functions to execute SQL queries and transaction
type SQLStore struct {
	*Queries
	db *sql.DB
}

// SQLTransaction 实现事务操作
type SQLTransaction struct {
	*Queries
	tx *sql.Tx
}

func NewStore(db *sql.DB) TxManager {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *SQLStore) Begin(ctx context.Context) (Tx, context.Context, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	x := &SQLTransaction{
		Queries: New(tx),
		tx:      tx,
	}
	ctx = withTx(ctx, x)
	return x, ctx, nil
}

func (t *SQLTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *SQLTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (s *SQLStore) Querier(ctx context.Context) Querier {
	if tx, ok := ctx.Value(txKey{}).(Tx); ok {
		return tx
	}
	return s
}

func withTx(ctx context.Context, tx Tx) context.Context {
	ctx = context.WithValue(ctx, txKey{}, tx)
	return ctx
}

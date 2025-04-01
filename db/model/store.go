package model

import (
	"context"
	"database/sql"
	"fmt"
)

type TxKey struct{} // 唯一的事务键

// TxManager 定义事务管理接口
type TxManager interface {
	Querier
	Begin(ctx context.Context) (Tx, error)
	ExecTx(ctx context.Context, fn func(Tx) error) error
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

func (s *SQLStore) Begin(ctx context.Context) (Tx, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &SQLTransaction{
		Queries: New(tx),
		tx:      tx,
	}, nil
}

func (t *SQLTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *SQLTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (m *SQLStore) ExecTx(ctx context.Context, fn func(Tx) error) error {
	tx, err := m.Begin(ctx)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

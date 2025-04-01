package model

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	BeginTx(ctx context.Context) error
	Commit() error
	Rollback() error
}

// SQLStore provides all functions to execute SQL queries and transaction
type SQLStore struct {
	*Queries
	db *sql.DB
	tx *sql.Tx
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SQLStore) BeginTx(ctx context.Context) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	store.Queries = New(tx)
	store.tx = tx
	return nil
}

func (store *SQLStore) Commit() error {
	if store.tx == nil {
		return fmt.Errorf("cannot commit: no transaction started")
	}

	err := store.tx.Commit()
	if err != nil {
		return err
	}

	store.Queries = New(store.db)
	store.tx = nil
	return nil
}

func (store *SQLStore) Rollback() error {
	if store.tx == nil {
		return fmt.Errorf("cannot rollback: no transaction started")
	}

	err := store.tx.Rollback()
	if err != nil {
		return err
	}

	store.Queries = New(store.db)
	store.tx = nil
	return nil
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v ,rb err:%v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

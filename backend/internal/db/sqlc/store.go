package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
	WithTx(tx *sql.Tx) Store
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func SetDBTX(db DBTX, store *Queries) {
	store.db = db
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// BeginTx starts a new transaction
func (store *SQLStore) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return store.db.BeginTx(ctx, nil)
}

// WithTx creates a new Store instance with the given transaction
func (store *SQLStore) WithTx(tx *sql.Tx) Store {
	return &SQLStore{
		Queries: New(tx),
		db:      store.db,
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("%w: %w, %w", err, err, rbErr)
		}
		return fmt.Errorf("%w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

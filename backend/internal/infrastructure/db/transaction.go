package db

import (
	"context"
	"database/sql"
)

type TxManager interface {
	RunInTx(ctx context.Context, f func(ctx context.Context) error) error
}

type txManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) TxManager {
	return &txManager{db: db}
}

func (m *txManager) RunInTx(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := f(ctx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

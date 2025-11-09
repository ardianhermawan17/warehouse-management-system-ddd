package sql

import (
	"context"
	"database/sql"
	"fmt"
)

// TransactionManager manages database transactions
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

// BeginTx starts a new transaction
func (tm *TransactionManager) BeginTx(ctx context.Context) (context.Context, error) {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return context.WithValue(ctx, "tx", tx), nil
}

// CommitTx commits the transaction
func (tm *TransactionManager) CommitTx(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return fmt.Errorf("no transaction in context")
	}

	return tx.Commit()
}

// RollbackTx rolls back the transaction
func (tm *TransactionManager) RollbackTx(ctx context.Context) error {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if !ok {
		return fmt.Errorf("no transaction in context")
	}

	return tx.Rollback()
}

// WithTx executes a function within a transaction
func (tm *TransactionManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txCtx := context.WithValue(ctx, "tx", tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// GetTx retrieves transaction from context
func GetTx(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if ok {
		return tx
	}
	return nil
}

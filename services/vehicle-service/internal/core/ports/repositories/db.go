package repositories

import (
	"context"
	"database/sql"
)

// DatabaseTransactioner is the minimal interface for executing SQL statements.
type DatabaseTransactioner interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Databaser wraps a full DB connection pool.
type Databaser interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	DatabaseTransactioner
	PingContext(ctx context.Context) error
	Close() error
}

// Transactioner wraps an active transaction.
type Transactioner interface {
	Commit() error
	Rollback() error
	DatabaseTransactioner
}

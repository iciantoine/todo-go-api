package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// DBTX is a database connection.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

// DB is a connection pool that we can create transactions from.
type DB interface {
	DBTX
	BeginTx(context.Context, *sql.TxOptions) (Tx, error)
}

// Tx is a transaction.
type Tx interface {
	DBTX
	Commit() error
	Rollback() error
}

// Used to make scanning consistent.
type scanner interface {
	Scan(args ...any) error
}

// list is meant to wrap QueryXXX that select a list of entities in the
// database. It returns a slice of these entities or nil if no rows found.
func list[T any](scan func(scanner) (T, error)) func(*sql.Rows, error) ([]T, error) {
	return func(rows *sql.Rows, err error) ([]T, error) {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, nil
		case err != nil:
			return nil, fmt.Errorf("could not execute query: %w", err)
		}
		defer rows.Close()

		var res []T
		for rows.Next() {
			val, err := scan(rows)
			if err != nil {
				return nil, fmt.Errorf("could not scan row: %w", err)
			}
			res = append(res, val)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		return res, nil
	}
}

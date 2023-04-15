package repository

import (
	"context"
	"database/sql"
	"portCaptureServer/app/entity"
)

type SavePortsRepository interface {
	StartTransaction() (Transaction, error)
	SavePort(ctx context.Context, transaction Transaction, port *entity.Port) error
}

// Transaction is added to help with testing
type Transaction interface {
	// QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Commit() error
	Rollback() error
}

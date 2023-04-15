package repository

import (
	"context"
	"portCaptureServer/app/entity"
)

type SavePortsRepository interface {
	StartTransaction() (Transaction, error)
	SavePort(ctx context.Context, transaction Transaction, port *entity.Port) error
}

// Transaction is added to help with testing
type Transaction interface {
	Commit() error
	Rollback() error
}

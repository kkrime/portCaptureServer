package sql

import (
	"context"
	"portCaptureServer/app/entity"
	"portCaptureServer/app/repository"

	"github.com/jmoiron/sqlx"
)

type savePortsRepository struct {
	db *sqlx.DB
}

func NewSavePortsRepository(db *sqlx.DB) repository.SavePortsRepository {

	return &savePortsRepository{
		db: db,
	}
}

func (spp *savePortsRepository) StartTransaction() (TransactionDB, error) {
	tx, err := spp.db.Begin()
	if err != nil {
		return nil, err
	}

	return NewTransactionDB(tx), nil
}

func (spp *savePortsRepository) SavePort(ctx context.Context, port *entity.Port) error {
	return savePort(ctx, spp.db, port)

}

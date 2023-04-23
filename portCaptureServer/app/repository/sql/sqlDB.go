package sql

import (
	"context"
	"portCaptureServer/app/entity"
	"portCaptureServer/app/repository"

	"github.com/jmoiron/sqlx"
)

type sqlDB struct {
	db *sqlx.DB
}

func NewSQLDB(db *sqlx.DB) repository.SavePortsRepository {

	return &sqlDB{
		db: db,
	}
}

func (spp *sqlDB) StartTransaction() (SQLTransactionDB, error) {
	tx, err := spp.db.Begin()
	if err != nil {
		return nil, err
	}

	return NewSQLTransactionDB(tx), nil
}

func (spp *sqlDB) SavePort(ctx context.Context, port *entity.Port) error {
	return savePort(ctx, spp.db, port)

}

package sql

import (
	"context"
	"database/sql"
	"portCaptureServer/app/entity"
)

type sqlTransactionDB struct {
	tx *sql.Tx
}

func NewSQLTransactionDB(tx *sql.Tx) SQLTransactionDB {
	return &sqlTransactionDB{
		tx: tx,
	}
}

func (tdb *sqlTransactionDB) SavePort(ctx context.Context, port *entity.Port) error {
	return savePort(ctx, tdb.tx, port)
}

func (tdb *sqlTransactionDB) Commit() error {
	return tdb.tx.Commit()
}

func (tdb *sqlTransactionDB) Rollback() error {
	return tdb.tx.Rollback()
}

package sql

import (
	"context"
	"database/sql"
	"portCaptureServer/app/entity"
)

type transactionDB struct {
	tx *sql.Tx
}

func NewTransactionDB(tx *sql.Tx) TransactionDB {
	return &transactionDB{
		tx: tx,
	}
}

func (tdb *transactionDB) SavePort(ctx context.Context, port *entity.Port) error {
	return savePort(ctx, tdb.tx, port)
}

func (tdb *transactionDB) Commit() error {
	return tdb.tx.Commit()
}

func (tdb *transactionDB) Rollback() error {
	return tdb.tx.Rollback()
}

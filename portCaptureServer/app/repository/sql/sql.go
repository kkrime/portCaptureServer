package sql

import "portCaptureServer/app/repository"

type DB interface {
	repository.SavePortsRepository
	StartTransaction() (SQLTransactionDB, error)
}

type SQLTransactionDB interface {
	repository.SavePortsRepository
	Commit() error
	Rollback() error
}

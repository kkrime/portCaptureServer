package sql

import "portCaptureServer/app/repository"

type DB interface {
	repository.SavePortsRepository
	StartTransaction() (TransactionDB, error)
}

type TransactionDB interface {
	repository.SavePortsRepository
	Commit() error
	Rollback() error
}

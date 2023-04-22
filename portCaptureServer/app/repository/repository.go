package repository

import (
	"context"
	"portCaptureServer/app/entity"
)

type SavePortsRepository interface {
	SavePort(ctx context.Context, port *entity.Port) error
}

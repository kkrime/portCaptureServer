package service

import "context"

type SendPortService interface {
	SendPort(ctx context.Context, portData *[]byte) error
}

package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *PortCaptureServer) HealthCheck(ctx context.Context, in *empty.Empty) (*empty.Empty, error) {

	return in, nil

}

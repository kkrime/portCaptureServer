package service

import (
	"context"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (sps *sendPortService) HealthCheck(ctx context.Context) error {

	_, err := sps.portCaptureServerClient.HealthCheck(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}

	return nil
}

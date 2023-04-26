package server

import (
	"context"
	"fmt"
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/service"

	pbMock "portCaptureServer/app/api/pb/mocks"
	serviceMock "portCaptureServer/app/service/mocks"
	sqlService "portCaptureServer/app/service/sql"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tj/assert"
)

type SavePortsTestSuite struct {
	suite.Suite
}

func TestSavePortsTestSuite(t *testing.T) {
	suite.Run(t, new(SavePortsTestSuite))
}

func (s *SavePortsTestSuite) TestSavePorts() {

	tests := []struct {
		name                     string
		portsStreamMock          pb.PortCaptureService_SavePortsServer
		savePortsServiceProvider service.SavePortsServiceProvider
		err                      error
		wrappedErr               error
	}{
		{
			name: "Happy Path",

			portsStreamMock: func() pb.PortCaptureService_SavePortsServer {
				portStreamMock := pbMock.NewPortCaptureService_SavePortsServer(s.T())

				portStreamMock.On(
					"SendAndClose",
					mock.Anything,
				).Return(
					nil,
				).
					Once()

				return portStreamMock
			}(),

			savePortsServiceProvider: func() service.SavePortsServiceProvider {
				savePortsServiceMock := serviceMock.NewSavePortsServiceProvider(s.T())

				savePortsServiceMock.On(
					"NewSavePortsInstance",
					mock.Anything,
					sqlService.SQLTransactionDB,
				).Return(
					func() service.SavePortsServiceInstance {
						savePortsServiceInstanceMock := serviceMock.NewSavePortsServiceInstance(s.T())

						savePortsServiceInstanceMock.On(
							"SavePorts",
							mock.Anything,
						).
							Return(
								nil,
							).
							Once()

						return savePortsServiceInstanceMock
					}(),
					nil,
				).
					Once()

				return savePortsServiceMock
			}(),
		},
		{
			name: "savePortsServiceProvider.NewSavePortsInstance() returns error",

			portsStreamMock: func() pb.PortCaptureService_SavePortsServer {
				portStreamMock := pbMock.NewPortCaptureService_SavePortsServer(s.T())

				portStreamMock.On(
					"SendAndClose",
					mock.Anything,
				).Return(
					nil,
				).
					Once()

				return portStreamMock
			}(),

			savePortsServiceProvider: func() service.SavePortsServiceProvider {
				savePortsServiceMock := serviceMock.NewSavePortsServiceProvider(s.T())

				savePortsServiceMock.On(
					"NewSavePortsInstance",
					mock.Anything,
					sqlService.SQLTransactionDB,
				).Return(
					nil,
					fmt.Errorf("NewSavePortsInstance() error"),
				).
					Once()

				return savePortsServiceMock
			}(),
		},
		{
			name: " savePortServiceInstance.SavePort() error",

			portsStreamMock: func() pb.PortCaptureService_SavePortsServer {
				portStreamMock := pbMock.NewPortCaptureService_SavePortsServer(s.T())

				portStreamMock.On(
					"SendAndClose",
					mock.Anything,
				).Return(
					nil,
				).
					Once()

				return portStreamMock
			}(),

			savePortsServiceProvider: func() service.SavePortsServiceProvider {
				savePortsServiceMock := serviceMock.NewSavePortsServiceProvider(s.T())

				savePortsServiceMock.On(
					"NewSavePortsInstance",
					mock.Anything,
					sqlService.SQLTransactionDB,
				).Return(
					func() service.SavePortsServiceInstance {
						savePortsServiceInstanceMock := serviceMock.NewSavePortsServiceInstance(s.T())

						savePortsServiceInstanceMock.On(
							"SavePorts",
							mock.Anything,
						).
							Return(
								fmt.Errorf("SavePort() error"),
							).
							Once()

						return savePortsServiceInstanceMock
					}(),
					nil,
				).
					Once()

				return savePortsServiceMock
			}(),
		},
		{
			name: "portsStream.SendAndClose() returns error",

			portsStreamMock: func() pb.PortCaptureService_SavePortsServer {
				portStreamMock := pbMock.NewPortCaptureService_SavePortsServer(s.T())

				portStreamMock.On(
					"SendAndClose",
					mock.Anything,
				).Return(
					fmt.Errorf("SendAndClose() error"),
				).
					Once()

				return portStreamMock
			}(),

			savePortsServiceProvider: func() service.SavePortsServiceProvider {
				savePortsServiceMock := serviceMock.NewSavePortsServiceProvider(s.T())

				savePortsServiceMock.On(
					"NewSavePortsInstance",
					mock.Anything,
					sqlService.SQLTransactionDB,
				).Return(
					nil,
					fmt.Errorf("NewSavePortsInstance() error"),
				).
					Once()

				return savePortsServiceMock
			}(),

			err: fmt.Errorf("SendAndClose() error"),
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t2 *testing.T) {
			server := NewPortCaptureServer(test.savePortsServiceProvider, context.Background())

			err := server.SavePorts(test.portsStreamMock)

			assert.Equal(s.T(), test.err, err)
		})
	}
}

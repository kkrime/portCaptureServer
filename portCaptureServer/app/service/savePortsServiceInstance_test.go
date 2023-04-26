package service

import (
	"context"
	"fmt"
	"io"
	"portCaptureServer/app/adapter"
	"portCaptureServer/app/entity"
	"portCaptureServer/app/logger"
	"portCaptureServer/app/repository"

	adapterMock "portCaptureServer/app/adapter/mocks"
	repoMock "portCaptureServer/app/repository/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tj/assert"
)

type SavePortsServiceInstanceTestSuite struct {
	suite.Suite
}

func TestSavePortsServiceInstanceTestSuite(t *testing.T) {
	suite.Run(t, new(SavePortsServiceInstanceTestSuite))
}

func (s *SavePortsServiceInstanceTestSuite) TestSendPortFinalize() {

	tests := []struct {
		name                    string
		portsStreamMock         adapter.PortsStream
		savePortsRepositoryMock repository.SavePortsRepository
		finalizeDB              func() error
		cancelDB                func() error
		err                     error
	}{
		{
			name: "Happy Path",

			portsStreamMock: func() adapter.PortsStream {
				portStreamMock := adapterMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&entity.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						&entity.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						io.EOF,
					)
				return portStreamMock
			}(),

			savePortsRepositoryMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
				).
					Return(
						nil,
					).
					Once()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
				).
					Return(
						nil,
					).
					Once()

				return repositoryMock
			}(),
		},
		{
			name: "SavePortsRepository.SavePort() returns error on first call",

			portsStreamMock: func() adapter.PortsStream {
				portStreamMock := adapterMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&entity.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						&entity.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						io.EOF,
					)
				return portStreamMock
			}(),

			savePortsRepositoryMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
				).
					Return(
						fmt.Errorf("SavePortsRepository.SavePort() error"),
					).
					Once()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
				).
					Return(
						nil,
					).
					// race condition
					Maybe()

				return repositoryMock
			}(),

			err: fmt.Errorf("SavePortsRepository.SavePort() error"),
		},
		{
			name: "SavePortsRepository.SavePort() returns error on second call",

			portsStreamMock: func() adapter.PortsStream {
				portStreamMock := adapterMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&entity.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						&entity.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						io.EOF,
					)
				return portStreamMock
			}(),

			savePortsRepositoryMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
				).
					Return(
						nil,
					).
					Once().
					// race condition
					Maybe()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
				).
					Return(
						fmt.Errorf("SavePortsRepository.SavePort() error"),
					).
					Once()

				return repositoryMock
			}(),

			err: fmt.Errorf("SavePortsRepository.SavePort() error"),
		},
		{
			name: "finalizeDB() fails",

			portsStreamMock: func() adapter.PortsStream {
				portStreamMock := adapterMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&entity.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						io.EOF,
					)
				return portStreamMock
			}(),

			savePortsRepositoryMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
				).
					Return(
						nil,
					).
					Once()

				return repositoryMock
			}(),

			finalizeDB: func() error {
				return fmt.Errorf("finalizeDB() error")
			},

			err: fmt.Errorf("finalizeDB() error"),
		},
		{
			name: "Recv() fails + cancelDB() fails",

			portsStreamMock: func() adapter.PortsStream {
				portStreamMock := adapterMock.NewPortsStream(s.T())

				portStreamMock.On(
					"Recv",
				).
					Return(
						&entity.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						fmt.Errorf("Recv() error"),
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						io.EOF,
					)

				return portStreamMock
			}(),

			savePortsRepositoryMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
				).
					Return(
						nil,
					).
					Once()

				return repositoryMock
			}(),

			cancelDB: func() error {
				return fmt.Errorf("cancelDB() error")
			},

			err: func() error {
				err1 := fmt.Errorf("Recv() error")
				err2 := fmt.Errorf("cancelDB() error")
				return fmt.Errorf("%v: %w", err1, err2)
			}(),
		},
		{
			name: "Recv() fails + cancelDB() fails",

			portsStreamMock: func() adapter.PortsStream {
				portStreamMock := adapterMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&entity.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						io.EOF,
					)

				return portStreamMock
			}(),

			savePortsRepositoryMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
				).
					Return(
						fmt.Errorf("repository.Recv() error"),
					).
					Once()

				return repositoryMock
			}(),

			cancelDB: func() error {
				return fmt.Errorf("cancelDB() error")
			},

			err: func() error {
				err1 := fmt.Errorf("repository.Recv() error")
				err2 := fmt.Errorf("cancelDB() error")
				return fmt.Errorf("%v: %w", err1, err2)
			}(),
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t2 *testing.T) {
			service := NewSavePortsServiceProvider(
				nil,
				2,
				logger.CreateNewLogger())

			savePortsServiceInstance := InitalizeNewSavePortsServiceInstance(
				context.TODO(),
				service.savePortsToDBChann,
				test.savePortsRepositoryMock,
				test.finalizeDB,
				test.cancelDB,
				logger.CreateNewLogger(),
			)

			err := savePortsServiceInstance.SavePorts(test.portsStreamMock)

			assert.Equal(s.T(), test.err, err)
		})
	}
}

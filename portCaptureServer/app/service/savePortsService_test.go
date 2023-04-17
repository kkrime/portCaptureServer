package service

import (
	"context"
	"fmt"
	"io"
	"portCaptureServer/app/api/pb"
	"portCaptureServer/app/logger"
	"portCaptureServer/app/repository"

	repoMock "portCaptureServer/app/repository/mocks"
	serviceMock "portCaptureServer/app/service/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/tj/assert"
)

type SavePortsServiceTestSuite struct {
	suite.Suite
}

func TestSavePortsServerTestSuite(t *testing.T) {
	suite.Run(t, new(SavePortsServiceTestSuite))
}

func (s *SavePortsServiceTestSuite) Testscan() {

	tests := []struct {
		name                     string
		savePortsTransactionMock repository.SavePortsRepository
		portsStreamMock          PortsStream
		err                      error
	}{
		{
			name: "Happy Path",

			savePortsTransactionMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"StartTransaction",
				).
					Return(
						func() repository.Transaction {
							TransactionMock := repoMock.NewTransaction(s.T())
							TransactionMock.On(
								"Commit",
							).
								Return(
									nil,
								).
								Once()

							return TransactionMock
						},
						nil,
					).
					Once()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						nil,
					).
					Once()

				return repositoryMock
			}(),

			portsStreamMock: func() PortsStream {
				portStreamMock := serviceMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&pb.Port{},
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
		},
		{
			name: "Recv() returns error on first call",

			savePortsTransactionMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"StartTransaction",
				).
					Return(
						func() repository.Transaction {
							TransactionMock := repoMock.NewTransaction(s.T())
							TransactionMock.On(
								"Rollback",
							).
								Return(
									nil,
								).
								Once()

							return TransactionMock
						},
						nil,
					).
					Once()

				return repositoryMock
			}(),

			portsStreamMock: func() PortsStream {
				portStreamMock := serviceMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						fmt.Errorf("Recv error"),
					).Once()

				return portStreamMock
			}(),

			err: fmt.Errorf("Recv error"),
		},
		{
			name: "Recv() returns error on second call",

			savePortsTransactionMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"StartTransaction",
				).
					Return(
						func() repository.Transaction {
							TransactionMock := repoMock.NewTransaction(s.T())
							TransactionMock.On(
								"Rollback",
							).
								Return(
									nil,
								).
								Once()

							return TransactionMock
						},
						nil,
					).
					Once()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						nil,
					).
					Once()

				return repositoryMock
			}(),

			portsStreamMock: func() PortsStream {
				portStreamMock := serviceMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&pb.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						fmt.Errorf("Recv error"),
					).Once()

				return portStreamMock
			}(),

			err: fmt.Errorf("Recv error"),
		},
		{
			name: "repository.SavePorts() fails on first call",

			savePortsTransactionMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"StartTransaction",
				).
					Return(
						func() repository.Transaction {
							TransactionMock := repoMock.NewTransaction(s.T())
							TransactionMock.On(
								"Rollback",
							).
								Return(
									nil,
								).
								Once()

							return TransactionMock
						},
						nil,
					).
					Once()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						fmt.Errorf("repository.SavePort() failed"),
					).
					Once()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						nil,
					).
					Once()

				return repositoryMock
			}(),

			portsStreamMock: func() PortsStream {
				portStreamMock := serviceMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&pb.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						&pb.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						io.EOF,
					).Once()

				return portStreamMock
			}(),

			err: fmt.Errorf("repository.SavePort() failed"),
		},
		{
			name: "repository.SavePorts() fails on second call",

			savePortsTransactionMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"StartTransaction",
				).
					Return(
						func() repository.Transaction {
							TransactionMock := repoMock.NewTransaction(s.T())
							TransactionMock.On(
								"Rollback",
							).
								Return(
									nil,
								).
								Once()

							return TransactionMock
						},
						nil,
					).
					Once()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
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
					mock.Anything,
				).
					Return(
						fmt.Errorf("repository.SavePort() failed"),
					).
					Once()

				return repositoryMock
			}(),

			portsStreamMock: func() PortsStream {
				portStreamMock := serviceMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&pb.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						&pb.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						io.EOF,
					).Once()

				return portStreamMock
			}(),

			err: fmt.Errorf("repository.SavePort() failed"),
		},
		{
			name: "repository.SavePorts() fails on first and second call",

			savePortsTransactionMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"StartTransaction",
				).
					Return(
						func() repository.Transaction {
							TransactionMock := repoMock.NewTransaction(s.T())
							TransactionMock.On(
								"Rollback",
							).
								Return(
									nil,
								).
								Once()

							return TransactionMock
						},
						nil,
					).
					Once()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						fmt.Errorf("repository.SavePort() failed"),
					).
					Once()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						fmt.Errorf("repository.SavePort() failed"),
					).
					Once()

				return repositoryMock
			}(),

			portsStreamMock: func() PortsStream {
				portStreamMock := serviceMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&pb.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						&pb.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						nil,
						io.EOF,
					).Once()

				return portStreamMock
			}(),

			err: fmt.Errorf("repository.SavePort() failed"),
		},
		{
			name: "Both Recv() returns error + repository.SavePorts() fails",

			savePortsTransactionMock: func() repository.SavePortsRepository {
				repositoryMock := repoMock.NewSavePortsRepository(s.T())

				repositoryMock.On(
					"StartTransaction",
				).
					Return(
						func() repository.Transaction {
							TransactionMock := repoMock.NewTransaction(s.T())
							TransactionMock.On(
								"Rollback",
							).
								Return(
									nil,
								).
								Once()

							return TransactionMock
						},
						nil,
					).
					Once()

				repositoryMock.On(
					"SavePort",
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).
					Return(
						fmt.Errorf("repository.SavePort() failed"),
					).
					Once()

				return repositoryMock
			}(),

			portsStreamMock: func() PortsStream {
				portStreamMock := serviceMock.NewPortsStream(s.T())
				portStreamMock.On(
					"Recv",
				).
					Return(
						&pb.Port{},
						nil,
					).Once()

				portStreamMock.On(
					"Recv",
				).
					Return(
						&pb.Port{},
						fmt.Errorf("Recv error"),
					).Once()

				return portStreamMock
			}(),

			err: fmt.Errorf("Recv error"),
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t2 *testing.T) {
			service := NewSavePortsService(test.savePortsTransactionMock, 2, logger.CreateNewLogger())
			err := service.SavePorts(context.Background(), test.portsStreamMock)

			assert.Equal(s.T(), test.err, err)
		})
	}
}

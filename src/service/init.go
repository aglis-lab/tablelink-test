package service

import (
	"tablelink/src/grpc"
	"tablelink/src/repository"
	"tablelink/transaction"
)

type UsersService interface {
	grpc.UsersServer
	// pb.EchoServer

	// BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error
}

type usersService struct {
	grpc.UnimplementedUsersServer

	authenticationService AuthenticationService

	userRepository      repository.UserRepository
	userRoleRepository  repository.UserRoleRepository
	roleRightRepository repository.RoleRightRepository
	gormTransaction     transaction.GormTransactionRepository
}

func NewUsersService(
	authenticationService AuthenticationService,
	userRepository repository.UserRepository,
	roleRightRepository repository.RoleRightRepository,
	userRoleRepository repository.UserRoleRepository,
	gormTransaction transaction.GormTransactionRepository,
) UsersService {
	return &usersService{
		authenticationService: authenticationService,
		userRepository:        userRepository,
		userRoleRepository:    userRoleRepository,
		roleRightRepository:   roleRightRepository,
		gormTransaction:       gormTransaction,
	}
}

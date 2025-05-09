package service

import (
	"tablelink/src/grpc"
	"tablelink/src/repository"
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
}

func NewUsersService(
	authenticationService AuthenticationService,
	userRepository repository.UserRepository,
	roleRightRepository repository.RoleRightRepository,
	userRoleRepository repository.UserRoleRepository,
) UsersService {
	return &usersService{
		authenticationService: authenticationService,
		userRepository:        userRepository,
		userRoleRepository:    userRoleRepository,
		roleRightRepository:   roleRightRepository,
	}
}

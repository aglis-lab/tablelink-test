package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"tablelink/src/app"
	"tablelink/src/repository"
	"tablelink/src/service"
	"tablelink/transaction"

	tablelink_grpc "tablelink/src/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type repositories struct {
	gormTransaction    transaction.GormTransactionRepository
	userRepository     repository.UserRepository
	rolesRepository    repository.RoleRightRepository
	userRoleRepository repository.UserRoleRepository
}

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Init app context
	if err := app.Init(ctx); err != nil {
		log.Panic(err)
	}

	// Init Router
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", app.Config().GrpcPort))
	if err != nil {
		log.Fatalln(err)
	}

	// Init grpc
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// Register Service into gRPC
	userService := newService()
	// grpcServer := grpc.NewServer(grpc.UnaryInterceptor(service.UnaryInterceptor))
	// pb.RegisterEchoServer(grpcServer, userService)
	tablelink_grpc.RegisterUsersServer(grpcServer, userService)

	fmt.Println("Listening gRPC", fmt.Sprintf("%s:%d", "tcp", app.Config().GrpcPort))
	if err := grpcServer.Serve(listen); err != nil {
		log.Println("error when running GRPC Server", err)
		log.Fatalln(err)
	}
}

func newService() service.UsersService {
	// Init Repository
	repo := newRepositories()

	// Service
	authenticationServie := service.NewAuthenticationService(repo.userRepository, repo.userRoleRepository, repo.rolesRepository)

	// Instanciate the service grpc
	return service.NewUsersService(
		authenticationServie,
		repo.userRepository,
		repo.rolesRepository,
		repo.userRoleRepository,
		repo.gormTransaction,
	)
}

func newRepositories() repositories {
	return repositories{
		gormTransaction:    transaction.NewGormTransactionRepository(app.GormDB()),
		userRepository:     repository.NewUserRepository(app.GormDB()),
		rolesRepository:    repository.NewRoleRightRepository(app.GormDB()),
		userRoleRepository: repository.NewUserRoleRepository(app.GormDB()),
	}
}

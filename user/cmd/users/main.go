package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mountolive/back-blog-go/user/grpc/transport"
	"github.com/mountolive/back-blog-go/user/pgstore"
	"github.com/mountolive/back-blog-go/user/usecase"
	"github.com/mountolive/back-blog-go/user/validation"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-shutdown
		fmt.Println("shutting down the server, users")
		cancel()
		log.Fatalf("users server shutting down...")
	}()
	dbUser := os.Getenv("USERS_DB_USER")
	dbPassword := os.Getenv("USERS_DB_PASS")
	dbName := os.Getenv("USERS_DB_NAME")
	dbPort := os.Getenv("USERS_DB_PORT")
	dbHost := os.Getenv("USERS_DB_HOST")
	dbUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)
	store, err := pgstore.NewUserPgStore(ctx, dbUrl)
	if err != nil {
		log.Fatalf("users db conn: %v", err)
	}
	validator, err := validation.NewValidator()
	if err != nil {
		log.Fatalf("users validator setup: %v", err)
	}
	repo := &usecase.UserRepository{
		Store:     store,
		Validator: validator,
	}
	_, err = repo.CreateUser(
		ctx,
		&usecase.CreateUserDto{
			Email:            os.Getenv("USERS_ADMIN_EMAIL"),
			Username:         os.Getenv("USERS_ADMIN_USERNAME"),
			Password:         os.Getenv("USERS_ADMIN_PASSWORD"),
			RepeatedPassword: os.Getenv("USERS_ADMIN_PASSWORD"),
			FirstName:        os.Getenv("USERS_ADMIN_FIRST_NAME"),
			LastName:         os.Getenv("USERS_ADMIN_LAST_NAME"),
		},
	)
	if err != nil {
		log.Printf("unable to create admin: %v\n", err)
	}
	gRPCServer := transport.NewGRPCServer(repo)
	baseServer := grpc.NewServer()
	// Same gRPC server will resolve all usecases
	transport.RegisterUserCheckerServer(baseServer, gRPCServer)
	transport.RegisterUserCreatorServer(baseServer, gRPCServer)
	transport.RegisterUserUpdaterServer(baseServer, gRPCServer)
	transport.RegisterPasswordChangerServer(baseServer, gRPCServer)
	transport.RegisterLoginServer(baseServer, gRPCServer)
	serverPort := os.Getenv("USERS_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", serverPort))
	if err != nil {
		log.Fatalf("users gRPC listener: %v", err)
	}
	fmt.Printf("users server starting at port: %s\n", serverPort)
	if err := baseServer.Serve(listener); err != nil {
		log.Fatalf("users gRPC server: %v", err)
	}
}

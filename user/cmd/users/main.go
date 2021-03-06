package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mountolive/back-blog-go/user/pgstore"
	"github.com/mountolive/back-blog-go/user/transport"
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
	gRPCServer := transport.NewGRPCServer(repo)
	baseServer := grpc.NewServer()
	transport.RegisterUserServer(baseServer, gRPCServer)
	serverPort := os.Getenv("USERS_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", serverPort))
	if err != nil {
		log.Fatalf("users gRPC listener: %v", err)
	}
	if err := baseServer.Serve(listener); err != nil {
		log.Fatalf("users gRPC server: %v", err)
	}
}

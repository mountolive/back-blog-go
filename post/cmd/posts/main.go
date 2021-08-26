package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/mountolive/back-blog-go/post/broker"
	"github.com/mountolive/back-blog-go/post/command"
	"github.com/mountolive/back-blog-go/post/eventbus"
	"github.com/mountolive/back-blog-go/post/httpx"
	"github.com/mountolive/back-blog-go/post/pgstore"
	"github.com/mountolive/back-blog-go/post/sanitizer"
	"github.com/mountolive/back-blog-go/post/usecase"
	"github.com/mountolive/back-blog-go/post/user"
	"github.com/mountolive/back-blog-go/post/user/transport"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-shutdown
		fmt.Println("shutting down the server, posts")
		cancel()
	}()
	dbUser := os.Getenv("POSTS_DB_USER")
	dbPassword := os.Getenv("POSTS_DB_PASS")
	dbName := os.Getenv("POSTS_DB_NAME")
	dbPort := os.Getenv("POSTS_DB_PORT")
	dbHost := os.Getenv("POSTS_DB_HOST")
	dbUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)
	store, err := pgstore.NewPostPgStore(ctx, dbUrl)
	if err != nil {
		log.Fatalf("posts db conn: %v", err)
	}
	gRPCHost := os.Getenv("POSTS_USERS_GRPC_HOST")
	gRPCPort := os.Getenv("POSTS_USERS_GRPC_PORT")
	gRPCConn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", gRPCHost, gRPCPort),
		// TODO Make gRPC connection to users' server to be secured
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("posts gRPC conn: %v", err)
	}
	client := transport.NewUserCheckerClient(gRPCConn)
	checker := user.NewGRPCUserChecker(client)
	repo := &usecase.PostRepository{
		Store:     store,
		Checker:   checker,
		Sanitizer: sanitizer.NewSanitizer(),
	}
	eventBus := eventbus.NewEventBus()
	eventBus.Register(command.CreatePostEventNameV1, command.NewCreatePost(repo))
	eventBus.Register(command.UpdatePostEventNameV1, command.NewUpdatePost(repo))
	// milliseconds
	pollingTime := 250
	port := os.Getenv("POSTS_NATS_PORT")
	natsPort, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("posts nats port parsing: %v", err)
	}
	natsConf := broker.NewNATSConfig(
		os.Getenv("POSTS_NATS_USER"),
		os.Getenv("POSTS_NATS_PASS"),
		os.Getenv("POSTS_NATS_SUBSCRIPTION_NAME"),
		os.Getenv("POSTS_NATS_DEADLETTER_NAME"),
		os.Getenv("POSTS_NATS_HOST"),
		uint16(natsPort),
		pollingTime,
	)
	natsBroker, err := broker.NewNATSBroker(eventBus, natsConf)
	if err != nil {
		log.Fatalf("posts nats broker: %v", err)
	}
	go func() {
		errChan := natsBroker.Process(ctx)
		for err := range errChan {
			fmt.Printf("posts nats process: %v\n", err)
		}
	}()
	httpServer := httpx.NewServer(repo)
	router := httpx.NewRouter()
	err = router.Add("GET /posts/([A-Za-z0-9]*)", httpServer.GetPost)
	if err != nil {
		log.Fatalf("posts router register, post by id: %v", err)
	}
	err = router.Add("GET /posts([\\?&_=A-Za-z0-9]*)?", httpServer.Filter)
	if err != nil {
		log.Fatalf("posts router register, post by tag and date: %v", err)
	}
	httpPort := os.Getenv("POSTS_HTTP_PORT")
	fmt.Printf("posts, starting http server at %s\n", httpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", httpPort), router); err != nil {
		log.Fatalf("posts http listen and serve: %v", err)
	}
}

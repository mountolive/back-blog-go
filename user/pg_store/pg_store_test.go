package pg_store

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/mountolive/back-blog-go/user/usecase"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

var store *PgStore

func TestMain(m *testing.M) {
	var err error
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Error starting docker: %s\n", err)
	}

	testPass := os.Getenv("POSTGRES_TEST_PASSWORD")
	testUser := os.Getenv("POSTGRES_TEST_USER")
	hostPort := "5433"
	containerPort := "5432"
	containerName := "blog_user_test"

	runOptions := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.1",
		Env: []string{
			fmt.Sprintf("POSTGRES_PASSWORD=%s", testPass),
			fmt.Sprintf("POSTGRES_USER=%s", testUser),
		},
		ExposedPorts: []string{containerPort},
		Name:         containerName,
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(containerPort): {
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			},
		},
	}
	container, err := pool.RunWithOptions(runOptions)

	ctx, cancel := context.WithCancel(context.Background())
	retryFunc := func() error {
		store, err = NewUserPgStore(ctx,
			fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s?sslmode=disable",
				testUser, testPass, hostPort, testUser),
		)
		fmt.Println(err)
		return err
	}
	if err := pool.Retry(retryFunc); err != nil {
		pool.RemoveContainerByName(containerName)
		log.Fatalf("An error occurred initializing the db: %s\n", err)
	}

	exitCode := m.Run()

	// Can't use defer with os.Exit()
	if err := pool.Purge(container); err != nil {
		pool.RemoveContainerByName(containerName)
		log.Fatalf("Error purging the container: %s\n", err)
	}

	pool.RemoveContainerByName(containerName)
	cancel()
	os.Exit(exitCode)
}

func TestPgStore(t *testing.T) {
	t.Run("Canary", func(t *testing.T) {
		var _ usecase.UserStore = &PgStore{}
		fmt.Println(store)
	})
}

package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/mountolive/back-blog-go/user/usecase"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"
)

var store *PgStore

func testMainWrapper(m *testing.M) int {
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
	defer func() {
		if err := pool.Purge(container); err != nil {
			pool.RemoveContainerByName(containerName)
			log.Fatalf("Error purging the container: %s\n", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	return m.Run()
}

func getCount(s *PgStore) int64 {
	var count int64
	row := s.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM users")
	row.Scan(&count)
	return count
}

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func TestPgStore(t *testing.T) {
	genericErr := "Got: %s \n Expected: %s"
	t.Run("Canary", func(t *testing.T) {
		var _ usecase.UserStore = &PgStore{}
	})

	t.Run("Create", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		data := &usecase.CreateUserDto{
			Email:     "abc@gmail.com",
			Password:  "test123456",
			Username:  "everlong",
			FirstName: "hello",
			LastName:  "breathout",
		}
		result, err := store.Create(ctx, data)
		require.True(t, err == nil, "An error was returned %s", err)
		require.True(t, result != nil, "No instance was returned from create")
		require.True(t, result.Id != "", "User's Id was not properly created")
		require.True(t, result.Email == data.Email, genericErr, result.Email, data.Email)
		require.True(t, result.Username == data.Username, genericErr, result.Username, data.Username)
		require.True(t, result.FirstName == data.FirstName,
			genericErr, result.FirstName, data.FirstName)
		require.True(t, result.LastName == data.LastName, genericErr, result.LastName, data.LastName)
		require.False(t, result.CreatedAt.IsZero(), "CreatedAt should have been set")
		require.False(t, result.UpdatedAt.IsZero(), "UpdatedAt should have been set")
		require.True(t, getCount(store) == 1, "Number of rows in db should be 1")
	})

	t.Run("Update", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		baseUser := &usecase.CreateUserDto{
			Email:     "test_update@gmail.com",
			Password:  "test123456",
			Username:  "test_update",
			FirstName: "test",
			LastName:  "test",
		}
		createResult, err := store.Create(ctx, baseUser)
		require.True(t, err == nil, "An error was returned on create: %s", err)
		require.True(t, createResult != nil, "No instance was returned from create")
		require.True(t, createResult.Id != "", "User's Id was not properly created")

		data := &usecase.UpdateUserDto{
			Email:     "somethingelse@gmail.com",
			Username:  "hello",
			FirstName: "this",
			LastName:  "is it",
		}
		rowCount := getCount(store)
		require.True(t, rowCount >= 1, "Number of rows in db should be 1 or more %d", rowCount)
		result, err := store.Update(ctx, createResult.Id, data)
		require.True(t, err == nil, "An error was returned on update: %s", err)
		require.True(t, result != nil, "No instance was returned from update")
		require.True(t, result.Email == data.Email, genericErr, result.Email, data.Email)
		require.True(t, result.Username == data.Username, genericErr, result.Username, data.Username)
		require.True(t, result.FirstName == data.FirstName,
			genericErr, result.FirstName, data.FirstName)
		require.True(t, result.LastName == data.LastName, genericErr, result.LastName, data.LastName)
		require.False(t, result.CreatedAt.IsZero(), "CreatedAt should have been set")
		require.False(t, result.UpdatedAt.IsZero(), "UpdatedAt should have been set")
		require.True(t, result.CreatedAt.Before(result.UpdatedAt), "UpdatedAt date was not updated")
		rowCount = getCount(store)
		require.True(t, rowCount >= 1, "Number of rows in db should be 1 or more %d", rowCount)
	})

	t.Run("UpdatePassword CheckIfCorrectPassword", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		baseUser := &usecase.CreateUserDto{
			Email:     "test_password@gmail.com",
			Password:  "test123456",
			Username:  "test_password",
			FirstName: "test",
			LastName:  "test",
		}
		createResult, err := store.Create(ctx, baseUser)
		require.True(t, err == nil, "An error was returned on create: %s", err)
		require.True(t, createResult != nil, "No instance was returned from create")
		require.True(t, createResult.Id != "", "User's Id was not properly created")

		rowCount := getCount(store)
		require.True(t, rowCount >= 1, "Number of rows in db should be 1 or more %d", rowCount)

		first := &usecase.ChangePasswordDto{
			Email:            "test_password@gmail.com",
			NewPassword:      "new123456",
			RepeatedPassword: "new123456",
			OldPassword:      "test123456",
		}
		err = store.UpdatePassword(ctx, first)
		require.True(t, err == nil, "An error was returned on update password, first: %s", err)

		firstChecker := &usecase.CheckUserAndPasswordDto{
			Email:    "test_password@gmail.com",
			Username: "test_password",
			Password: "new123456",
		}
		err = store.CheckIfCorrectPassword(ctx, firstChecker)
		require.True(t, err == nil, "An error was returned on first check for password: %s", err)

		second := &usecase.ChangePasswordDto{
			Username:         "test_password",
			NewPassword:      "evenNewer123456",
			RepeatedPassword: "evenNewer123456",
			OldPassword:      "new123456",
		}
		err = store.UpdatePassword(ctx, second)
		require.True(t, err == nil, "An error was returned on update password, second: %s", err)

		secondChecker := &usecase.CheckUserAndPasswordDto{
			Email:    "test_password@gmail.com",
			Username: "test_password",
			Password: "evenNewer123456",
		}
		err = store.CheckIfCorrectPassword(ctx, secondChecker)
		require.True(t, err == nil, "An error was returned on second check for password: %s", err)
	})

	t.Run("ReadOne", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		baseUser := &usecase.CreateUserDto{
			Email:     "test_read@gmail.com",
			Password:  "test123456",
			Username:  "test_read",
			FirstName: "test",
			LastName:  "test",
		}
		created, err := store.Create(ctx, baseUser)
		require.True(t, err == nil, "An error was returned on create: %s", err)
		require.True(t, created != nil, "No instance was returned from create")
		require.True(t, created.Id != "", "User's Id was not properly created")

		rowCount := getCount(store)
		require.True(t, rowCount >= 1, "Number of rows in db should be 1 or more %d", rowCount)

		firstChecker := &usecase.ByUsernameOrEmail{
			Email: "test_read@gmail.com",
		}
		result, err := store.ReadOne(ctx, firstChecker)
		require.True(t, err == nil, "An error was returned on check for password: %s", err)
		require.True(t, result != nil, "No instance was returned from read")
		require.True(t, created.Id == result.Id, "User's Id was not properly created")
		require.True(t, result.Email == created.Email, genericErr, result.Email, created.Email)
		require.True(t, result.Username == created.Username, genericErr,
			result.Username, created.Username)

		secondChecker := &usecase.ByUsernameOrEmail{
			Username: "test_read",
		}
		result, err = store.ReadOne(ctx, secondChecker)
		require.True(t, err == nil, "An error was returned on check for password: %s", err)
		require.True(t, created.Id == result.Id, "User's Id was not properly created")
		require.True(t, result.Email == created.Email, genericErr, result.Email, created.Email)
		require.True(t, result.Username == created.Username, genericErr,
			result.Username, created.Username)
	})
}

package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/mountolive/back-blog-go/post/usecase"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/require"
)

const genericErr = "\nGot: %s \n Expected: %s\n"

var store *PgStore

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func TestPgStore(t *testing.T) {
	t.Run("Canary", func(t *testing.T) {
		var _ usecase.PostStore = &PgStore{}
	})

	t.Run("Create", func(t *testing.T) {
		post := &usecase.CreatePostDto{
			Creator: "theUser",
			Content: "anything",
			Tags:    []string{"tag1", "tag2"},
		}
		result := createPost(t, post)
		for _, tag := range post.Tags {
			checkPostsByTag(t, result, tag)
		}
	})

	t.Run("Update", func(t *testing.T) {
		post := &usecase.CreatePostDto{
			Creator: "sonic",
			Content: "Incinerate",
			Tags:    []string{"tag1", "tag2"},
		}
		result := createPost(t, post)

		updatedPost := &usecase.UpdatePostDto{
			Id:      result.Id,
			Content: "Bull in the heather",
			Tags:    []string{"tag3"},
		}
		updated, err := store.Update(context.Background(),
			updatedPost)
		require.True(t, err == nil, "Error was returned. Update")
		require.True(t, updated != nil, "No entity returned, Update")
		require.True(t, updated.Id == result.Id, "Ids not matching after update")
		require.True(t, updated.Content == updatedPost.Content, "Content not updated")
		for _, tag := range updatedPost.Tags {
			checkPostsByTag(t, result, tag)
		}
	})
}

func createPost(t *testing.T, post *usecase.CreatePostDto) *usecase.PostDto {
	result, err := store.Create(context.Background(), post)
	require.True(t, err == nil, "An error was returned. Not expected, Create")
	require.True(t, result != nil, "No entity was returned from Create")
	require.True(t, result.Id != "", "Id was empty, error creating the Post")
	require.True(t, result.Creator == post.Creator, genericErr,
		result.Creator, post.Creator)
	require.True(t, result.Content == post.Content, genericErr,
		result.Content, post.Content)
	return result
}

func checkPostsByTag(t *testing.T, result *usecase.PostDto, tag string) {
	filter := &usecase.GeneralFilter{PageSize: 1}
	filter.Tag = tag
	filteredPosts, err := store.Filter(context.Background(), filter)
	found := len(filteredPosts)
	require.True(t, err == nil,
		"An error was returned. Not expected, Create's Filter")
	require.True(t, found == 1, genericErr, found, 1)
	require.True(t, filteredPosts[0].Id == result.Id,
		"Created post doesn't match with found post")
}

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
	containerName := "blog_post_test"

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
		store, err = NewPostPgStore(ctx,
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

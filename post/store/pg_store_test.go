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
			Title:   "knows all",
			Content: "anything",
			Tags:    []string{"tag1", "tag2"},
		}
		createPost(t, post)
	})

	t.Run("Update", func(t *testing.T) {
		post := &usecase.CreatePostDto{
			Creator: "sonic",
			Title:   "youth",
			Content: "Incinerate",
			Tags:    []string{"tag3", "tag4"},
		}
		result := createPost(t, post)

		updatedPost := &usecase.UpdatePostDto{
			Id:      result.Id,
			Content: "Bull in the heather",
			Title:   "Playing bass like Kim Gordon",
			Tags:    []string{"tag4"},
		}
		updated, err := store.Update(context.Background(),
			updatedPost)
		require.True(t, err == nil, "Error was returned. Update")
		require.True(t, updated != nil, "No entity returned, Update")
		require.True(t, updated.Id == result.Id, "Ids not matching after update")
		require.True(t, updated.Content == updatedPost.Content, "Content not updated")
		require.True(t, updated.Title == updatedPost.Title, "Title not updated")
		for _, tag := range updatedPost.Tags {
			checkPostsByTag(t, result, tag, 1)
		}
	})

	t.Run("Filter", func(t *testing.T) {
		posts := []*usecase.CreatePostDto{
			{
				Creator: "first",
				Title:   "firstT",
				Content: "hello",
				Tags:    []string{"tag5"},
			},
			{
				Creator: "second",
				Title:   "secondT",
				Content: "there",
				Tags:    []string{"tag6"},
			},
			{
				Creator: "third",
				Title:   "thirdT",
				Content: "nope",
				Tags:    []string{"tag7"},
			},
		}
		createdPosts := []*usecase.PostDto{}
		for _, newPost := range posts {
			createdPosts = append(createdPosts, createPost(t, newPost))
		}

		applyAndCheckFilter := func(filter *usecase.GeneralFilter) {
			list, err := store.Filter(context.Background(), filter)
			expectedLen := 2
			require.True(t, err == nil, "Error while filtering posts %s", err)
			require.True(t, list != nil, "Nil pointer for filtered posts slice")
			actualLen := len(list)
			require.True(t, actualLen == expectedLen, genericErr,
				actualLen, expectedLen)
			require.True(t, list[0].Creator != list[1].Creator,
				"Creators should be different")
			require.True(t, list[0].Content != list[1].Content,
				"Content should be different")
			includedContent := func(content string) bool {
				return content == posts[0].Content || content == posts[1].Content
			}
			includedCreator := func(creator string) bool {
				return creator == posts[0].Creator || creator == posts[1].Creator
			}
			for _, result := range list {
				require.True(t, includedContent(result.Content), "Incorrect content %s: %s",
					result.Id, result.Content)
				require.True(t, includedCreator(result.Creator), "Incorrect creator %s: %s",
					result.Id, result.Creator)
			}
		}

		postLength := len(posts)
		tagFilter := &usecase.GeneralFilter{PageSize: postLength}
		tagFilter.Tag = "tag1"
		applyAndCheckFilter(tagFilter)

		dateFilter := &usecase.GeneralFilter{PageSize: postLength}
		dateFilter.From = createdPosts[0].CreatedAt
		dateFilter.To = createdPosts[1].CreatedAt
		applyAndCheckFilter(dateFilter)
	})

	t.Run("ReadOne", func(t *testing.T) {
		post := &usecase.CreatePostDto{
			Creator: "melvins",
			Title:   "Buzzo",
			Content: "the bit",
			Tags:    []string{"stag"},
		}
		result := createPost(t, post)

		found := store.ReadOne(context.Background(), result.Id)
		require.True(t, found != nil, "Post not found by the passed Id. ReadOne")
		require.True(t, found.Content == post.Content, genericErr,
			found.Content, post.Content)
		require.True(t, found.Creator == post.Creator, genericErr,
			found.Creator, post.Creator)
	})
}

func createPost(t *testing.T, post *usecase.CreatePostDto) *usecase.PostDto {
	result, err := store.Create(context.Background(), post)
	t.Log(result)
	require.True(t, err == nil, "An error was returned. Not expected: %s, Create", err)
	require.True(t, result != nil, "No entity was returned from Create")
	require.True(t, result.Id != "", "Id was empty, error creating the Post")
	require.True(t, result.Creator == post.Creator, genericErr,
		result.Creator, post.Creator)
	require.True(t, result.Content == post.Content, genericErr,
		result.Content, post.Content)
	return result
}

func checkPostsByTag(t *testing.T, result *usecase.PostDto,
	tag string, pageSize int) {
	filter := &usecase.GeneralFilter{PageSize: pageSize}
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

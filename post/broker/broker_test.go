package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/mountolive/back-blog-go/post/eventbus"
	"github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func TestNATSConfig(t *testing.T) {
	t.Run("URL", func(t *testing.T) {
		t.Parallel()

		t.Run("Non-zero values", func(t *testing.T) {
			user := "testUsername"
			pass := "testPassword"
			port := uint16(16661)
			host := "localhost"
			expectedURL := fmt.Sprintf("nats://%s:%s@%s:%d", user, pass, host, port)
			actualURL := NewNATSConfig(user, pass, host, port).URL()
			require.Equal(t, expectedURL, actualURL)
		})

		t.Run("Zero values", func(t *testing.T) {
			expectedURL := nats.DefaultURL
			actualURL := NewNATSConfig("", "", "", 0).URL()
			require.Equal(t, expectedURL, actualURL)
			require.Equal(t, expectedURL, DefaultNATSConfig().URL())
		})
	})
}

var _ eventbus.Event = mockEvent{}

type mockEvent struct{}

func (mockEvent) Name() string { return "ready for my close-up" }

func (mockEvent) Params() eventbus.Params { return nil }

var _ EventBus = mockNonErroredEventBus{}

type mockNonErroredEventBus struct{}

func (mockNonErroredEventBus) Resolve(context.Context, eventbus.Event) error {
	return nil
}

var _ EventBus = mockErroredEventBus{}

type mockErroredEventBus struct {
	err error
}

func (m mockErroredEventBus) Resolve(context.Context, eventbus.Event) error {
	return m.err
}

func TestNATSBroker(t *testing.T) {
	t.Run("NewNATSBroker", func(t *testing.T) {
		t.Parallel()
		erroredBus := mockErroredEventBus{}
		notErroredBus := mockNonErroredEventBus{}

		t.Run("Connection error", func(t *testing.T) {
			conf := NewNATSConfig("badU", "badP", "badH", uint16(3333))
			_, err := NewNATSBroker(notErroredBus, conf)
			require.Error(t, err)
			require.True(t, errors.Is(err, ErrNATSServerConnection))
		})

		t.Run("Correct initialization", func(t *testing.T) {
			conf := DefaultNATSConfig()
			broker, err := NewNATSBroker(notErroredBus, conf)
			require.NoError(t, err)
			require.NotNil(t, broker)
		})
	})

	t.Run("ProcessMessages", func(t *testing.T) {
	})
}

func testMainWrapper(m *testing.M) int {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("error starting docker: %s\n", err)
	}

	containerName := "blog_post_nats"
	basePort := "4222"
	runOptions := &dockertest.RunOptions{
		Repository:   "nats",
		Tag:          "2.1",
		Name:         containerName,
		ExposedPorts: []string{basePort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(basePort): {
				{
					HostIP:   "0.0.0.0",
					HostPort: basePort,
				},
			},
		},
	}

	container, err := pool.RunWithOptions(runOptions)
	if err != nil {
		log.Fatalf("error occurred setting the container: %s", err)
	}
	defer func() {
		if err := pool.Purge(container); err != nil {
			errRemove := pool.RemoveContainerByName(containerName)
			if errRemove != nil {
				log.Fatalf("error removing the container: %s\n", errRemove)
			}
			_ = pool.RemoveContainerByName(containerName)
			log.Fatalf("error purging the container: %s\n", err)
		}
	}()

	retryFunc := func() error {
		conn, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			return err
		}
		defer conn.Close()
		return nil
	}
	if err := pool.Retry(retryFunc); err != nil {
		err2 := pool.RemoveContainerByName(containerName)
		if err2 != nil {
			log.Fatalf("error removing the container: %s\n", err2)
		}
		log.Fatalf("error occurred initializing the nats server: %s\n", err)
	}

	return m.Run()
}

package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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
			actualURL := NewNATSConfig(user, pass, "any", host, port).URL()
			require.Equal(t, expectedURL, actualURL)
		})

		t.Run("Zero values", func(t *testing.T) {
			expectedURL := nats.DefaultURL
			actualURL := NewNATSConfig("", "", "blaaa", "", 0).URL()
			require.Equal(t, expectedURL, actualURL)
			require.Equal(t, expectedURL, DefaultNATSConfig("something").URL())
		})
	})
}

var _ eventbus.Event = mockEvent{}

type mockEvent struct{}

func (mockEvent) Data() []byte {
	return []byte(
		`{"event_name": "hey", "data": "bla"}`,
	)
}

var _ EventBus = &mockNonErroredEventBus{}

type mockNonErroredEventBus struct {
	timesCalled int
}

func (m *mockNonErroredEventBus) Resolve(context.Context, eventbus.Event) error {
	m.timesCalled++
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
	mockErr := errors.New("I exploded")
	erroredBus := mockErroredEventBus{mockErr}
	notErroredBus := &mockNonErroredEventBus{}

	t.Run("NewNATSBroker", func(t *testing.T) {
		t.Parallel()

		t.Run("Connection error", func(t *testing.T) {
			conf := NewNATSConfig("badU", "badP", "any", "badH", uint16(3333))
			_, err := NewNATSBroker(notErroredBus, conf)
			require.Error(t, err)
			require.True(t, errors.Is(err, ErrNATSServerConnection))
		})

		t.Run("Correct initialization", func(t *testing.T) {
			conf := DefaultNATSConfig("some")
			broker, err := NewNATSBroker(notErroredBus, conf)
			require.NoError(t, err)
			require.NotNil(t, broker)
			// Avoid connection leakage
			broker.CloseConnection()
		})
	})

	t.Run("StartSubscription", func(t *testing.T) {
		t.Parallel()

		initializeBroker := func() *NATSBroker {
			conf := DefaultNATSConfig("bla")
			broker, err := NewNATSBroker(notErroredBus, conf)
			require.NoError(t, err)
			require.NotNil(t, broker)
			return broker
		}

		t.Run("Error starting subscription", func(t *testing.T) {
			broker := initializeBroker()
			broker.CloseConnection()
			err := broker.StartSubscription()
			require.True(t, errors.Is(err, ErrNATSubscription))
		})

		t.Run("Correct initialization", func(t *testing.T) {
			broker := initializeBroker()
			err := broker.StartSubscription()
			require.NoError(t, err)
		})
	})

	t.Run("ProcessMessages", func(t *testing.T) {
		t.Parallel()

		basicMsg := "some message %d"
		conf := DefaultNATSConfig("bla")
		producerConn, err := nats.Connect(conf.URL())
		require.NoError(t, err)

		produceMsg := func(n int) error {
			if producerConn.IsClosed() {
				return errors.New("connection closed")
			}
			return producerConn.Publish(
				conf.subscriptionName,
				[]byte(fmt.Sprintf(basicMsg, n)),
			)
		}
		timeout := 2 * time.Second

		brokerWithInitializedSubscription := func(bus EventBus) *NATSBroker {
			broker, err := NewNATSBroker(bus, conf)
			require.NoError(t, err)
			require.NotNil(t, broker)
			err = broker.StartSubscription()
			require.NoError(t, err)
			return broker
		}

		t.Run("Error from EventBus", func(t *testing.T) {
			broker := brokerWithInitializedSubscription(erroredBus)
			err = produceMsg(0)
			require.NoError(t, err)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			errCount := 0
			for err := range broker.Process(ctx) {
				fmt.Println("Error from EventBus:", err)
				errCount++
			}
			require.Equal(t, 2, errCount)
		})

		t.Run("Message consumption", func(t *testing.T) {
			broker := brokerWithInitializedSubscription(notErroredBus)
			var err error
			for i := 1; i < 11; i++ {
				err = produceMsg(i)
				require.NoError(t, err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			errCount := 0
			for err := range broker.Process(ctx) {
				fmt.Println("Message consumption", err)
				errCount++
			}
			require.Equal(t, 10, notErroredBus.timesCalled)
			require.Equal(t, 1, errCount)
		})
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

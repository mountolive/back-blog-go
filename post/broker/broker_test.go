package broker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
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
			actualURL := NewNATSConfig(
				user, pass, "any", "dead.any", host, port, 1,
			).URL()
			require.Equal(t, expectedURL, actualURL)
		})

		t.Run("Zero values", func(t *testing.T) {
			expectedURL := nats.DefaultURL
			actualURL := NewNATSConfig("", "", "blaaa", "dead.blaaa", "", 0, 1).URL()
			require.Equal(t, expectedURL, actualURL)
			require.Equal(t, expectedURL, DefaultNATSConfig("something").URL())
		})
	})
}

var _ eventbus.Event = mockEvent{}

type mockEvent struct {
	data string
}

func (m mockEvent) Data() []byte {
	return []byte(m.data)
}

var _ EventBus = &mockNonErroredEventBus{}

type mockNonErroredEventBus struct {
	timesCalled int
	resolveFunc func(context.Context, eventbus.Event) error
}

const testingMsg = "some_message_%d"

func (m *mockNonErroredEventBus) Resolve(ctx context.Context, ev eventbus.Event) error {
	if m.resolveFunc == nil {
		panic("resolveFunc has no implementation")
	}
	return m.resolveFunc(ctx, ev)
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
	notErroredBus.resolveFunc = func(ctx context.Context, ev eventbus.Event) error {
		require.Equal(
			t,
			fmt.Sprintf(testingMsg, notErroredBus.timesCalled),
			string(ev.Data()),
		)
		notErroredBus.timesCalled += 1
		return nil
	}

	t.Run("NewNATSBroker", func(t *testing.T) {
		t.Run("Connection error", func(t *testing.T) {
			conf := NewNATSConfig(
				"badU", "badP", "any", "dead.any", "badH", uint16(3333), 1,
			)
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

	t.Run("ProcessMessages", func(t *testing.T) {
		conf := DefaultNATSConfig("bla")
		producerConn, err := nats.Connect(conf.URL())
		require.NoError(t, err)

		produceMsg := func(msgNum int) error {
			if producerConn.IsClosed() {
				return errors.New("connection closed")
			}
			return producerConn.Publish(
				conf.subscriptionName,
				[]byte(fmt.Sprintf(testingMsg, msgNum)),
			)
		}
		timeout := 5 * time.Second

		brokerWithInitializedSubscription := func(bus EventBus) *NATSBroker {
			broker, err := NewNATSBroker(bus, conf)
			require.NoError(t, err)
			require.NotNil(t, broker)
			return broker
		}

		t.Run("Error from EventBus", func(t *testing.T) {
			broker := brokerWithInitializedSubscription(erroredBus)
			defer broker.CloseConnection()
			err = produceMsg(rand.Intn(100))
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
			defer broker.CloseConnection()
			var err error
			for i := 1; i < 11; i++ {
				err = produceMsg(i)
				require.NoError(t, err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			errCount := 0
			for err := range broker.Process(ctx) {
				fmt.Println("Message consumption:", err)
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
	inspectPort := "8222"
	runOptions := &dockertest.RunOptions{
		Repository:   "nats",
		Tag:          "2.1",
		Name:         containerName,
		ExposedPorts: []string{basePort, inspectPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(basePort): {
				{
					HostIP:   "0.0.0.0",
					HostPort: basePort,
				},
			},
			docker.Port(inspectPort): {
				{
					HostIP:   "0.0.0.0",
					HostPort: inspectPort,
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

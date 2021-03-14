package broker

import (
	"log"
	"os"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("error starting docker: %s\n", err)
	}

	hostIP := "0.0.0.0"
	containerName := "blog_post_nats"

	basePort := "4222"
	portBindings := make(map[docker.Port][]docker.PortBinding)
	for _, port := range []string{basePort, "6222", "8222"} {
		portBindings[docker.Port(port)] = []docker.PortBinding{
			{
				HostIP:   hostIP,
				HostPort: port,
			},
		}
	}

	runOptions := &dockertest.RunOptions{
		Repository:   "nats",
		Tag:          "2.1",
		Name:         containerName,
		PortBindings: portBindings,
	}

	container, err := pool.RunWithOptions(runOptions)
	if err != nil {
		log.Fatalf("error occurred setting the container: %s", err)
	}
	defer func() {
		if err := pool.Purge(container); err != nil {
			err2 := pool.RemoveContainerByName(containerName)
			if err2 != nil {
				log.Fatalf("error removing the container: %s\n", err2)
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

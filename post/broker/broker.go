package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/mountolive/back-blog-go/post/eventbus"
	"github.com/nats-io/nats.go"
)

var ErrNATSServerConnection = errors.New("nats server error connection")

// EventBus is the needed functionality from the corresponding Broker
type EventBus interface {
	Resolve(context.Context, eventbus.Event) error
}

// NATSBroker implementation of a Broker on top of NATS (https://nats.io/)
type NATSBroker struct {
	bus EventBus
}

// NATSConfig is the basic configuration for a NATS connection
type NATSConfig struct {
	port uint16
	user,
	pass,
	host string
	opts []nats.Options
}

// NewNATSConfig is a standard constructor
func NewNATSConfig(
	user, pass, host string,
	port uint16,
	opts ...nats.Options) NATSConfig {
	// TODO Implement
	return NATSConfig{}
}

func (n NATSConfig) URL() string {
	// TODO Implement
	return ""
}

// DefaultNATSConfig returns a standard configuration for a barebones NATS server
func DefaultNATSConfig() NATSConfig {
	return NATSConfig{
		user: "",
		pass: "",
		host: "",
		port: 0,
	}
}

// NewNATSBroker is a standard constructor
func NewNATSBroker(bus EventBus, conf NATSConfig) (*NATSBroker, error) {
	// TODO Implement
	return nil, nil
}

// ProcessMessages starts consuming messages from a given subscription
func ProcessMessages(context.Context) error {
	// TODO Implement
	return nil
}

func wrapError(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

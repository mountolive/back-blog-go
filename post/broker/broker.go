package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/mountolive/back-blog-go/post/eventbus"
	"github.com/nats-io/nats.go"
)

// ErrNATSServerConnection indicates that there was an error when connection to
// the NATS server
var ErrNATSServerConnection = errors.New("nats server error connection")

// EventBus is the needed functionality from the corresponding Broker
type EventBus interface {
	Resolve(context.Context, eventbus.Event) error
}

// NATSBroker implementation of a Broker on top of NATS (https://nats.io/)
type NATSBroker struct {
	bus  EventBus
	conn *nats.Conn
}

// NATSConfig is the basic configuration for a NATS connection
type NATSConfig struct {
	port uint16
	user,
	pass,
	host string
	opts []nats.Option
}

// NewNATSConfig is a standard constructor
func NewNATSConfig(usr, pwd, h string, p uint16, opts ...nats.Option) NATSConfig {
	return NATSConfig{
		user: usr,
		pass: pwd,
		host: h,
		port: p,
		opts: opts,
	}
}

// URL returns the necessary URL to connect to a given NATS server
func (n NATSConfig) URL() string {
	port := n.port
	if port == 0 {
		port = 4222
	}
	if n.user == "" {
		if n.host == "" {
			return fmt.Sprintf("nats://127.0.0.1:%d", port)
		}
		return fmt.Sprintf("nats://%s:%d", n.host, port)
	}
	return fmt.Sprintf("nats://%s:%s@%s:%d", n.user, n.pass, n.host, port)
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
	conn, err := nats.Connect(conf.URL(), conf.opts...)
	if err != nil {
		return nil, wrapError(ErrNATSServerConnection, err.Error())
	}
	return &NATSBroker{bus: bus, conn: conn}, nil
}

// CloseConnection closes the underlying NATS connection;
// meant for cleanUp purposes
func (n *NATSBroker) CloseConnection() {
	n.conn.Close()
}

// ProcessMessages starts consuming messages from a given subscription
func (n *NATSBroker) ProcessMessages(context.Context) error {
	// TODO Implement
	return nil
}

func wrapError(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/mountolive/back-blog-go/post/eventbus"
	"github.com/nats-io/nats.go"
)

var (
	// ErrNATSServerConnection indicates that there was an error when connection to
	// the NATS server
	ErrNATSServerConnection = errors.New("NATS server error connection")
	// ErrNATSubscription indicates than an error occurred when triggering subscription
	ErrNATSubscription = errors.New("NATS subscription failed")
	// ErrEventBus
	ErrEventBus = errors.New("event bus error")
)

// EventBus is the needed functionality from the corresponding Broker
type EventBus interface {
	Resolve(context.Context, eventbus.Event) error
}

// NATSBroker implementation of a Broker on top of NATS (https://nats.io/)
type NATSBroker struct {
	bus          EventBus
	conn         *nats.Conn
	conf         NATSConfig
	messagesChan chan *nats.Msg
}

// NATSConfig is the basic configuration for a NATS connection
type NATSConfig struct {
	port uint16
	user,
	pass,
	subscriptionName,
	host string
	opts []nats.Option
}

// NewNATSConfig is a standard constructor
func NewNATSConfig(
	usr, pwd, subsName, h string, p uint16, opts ...nats.Option,
) NATSConfig {
	return NATSConfig{
		user:             usr,
		pass:             pwd,
		subscriptionName: subsName,
		host:             h,
		port:             p,
		opts:             opts,
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
// for the passed subscription's name
func DefaultNATSConfig(subsName string) NATSConfig {
	return NATSConfig{
		user:             "",
		pass:             "",
		subscriptionName: subsName,
		host:             "",
		port:             0,
	}
}

// NewNATSBroker is a standard constructor
func NewNATSBroker(bus EventBus, conf NATSConfig) (*NATSBroker, error) {
	conn, err := nats.Connect(conf.URL(), conf.opts...)
	if err != nil {
		return nil, wrapError(ErrNATSServerConnection, err.Error())
	}
	broker := &NATSBroker{
		bus:          bus,
		conn:         conn,
		conf:         conf,
		messagesChan: make(chan *nats.Msg),
	}
	return broker, nil
}

// StartSubscription starts the associated subscription
func (n *NATSBroker) StartSubscription() error {
	_, err := n.conn.ChanSubscribe(n.conf.subscriptionName, n.messagesChan)
	if err != nil {
		return wrapError(ErrNATSubscription, err.Error())
	}
	return nil
}

// CloseConnection closes the underlying NATS connection;
// meant for cleanUp purposes
func (n *NATSBroker) CloseConnection() {
	n.conn.Close()
}

// ProcessMessages starts cosuming messages from a given subscription
func (n *NATSBroker) Process(context.Context) error {
	// TODO Implement
	return nil
}

func wrapError(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

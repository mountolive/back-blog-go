package broker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mountolive/back-blog-go/post/eventbus"
	"github.com/nats-io/nats.go"
)

var (
	// ErrNATSServerConnection indicates that there was an error when connection to
	// the NATS server
	ErrNATSServerConnection = errors.New("NATS server error connection")
	// ErrNATSubscription indicates than an error occurred when triggering subscription
	ErrNATSubscription = errors.New("NATS subscription failed")
	// ErrEventBus is self-described
	ErrEventBus = errors.New("event bus error")
	// ErrContextCanceled is self-described
	ErrContextCanceled = errors.New("NATS broker context canceled")
	// ErrFlushSubscription is self-described
	ErrFlushSubscription = errors.New("NATS in flushing subscription (roundtrip)")
	// ErrDeadLetterPublish is self-described
	ErrDeadLetterPublish = errors.New("NATS publish to dead letter failed")
)

// EventBus is the needed functionality from the corresponding Broker
type EventBus interface {
	Resolve(context.Context, eventbus.Event) error
}

// NATSConfig is the basic configuration for a NATS connection
type NATSConfig struct {
	port uint16
	user,
	pass,
	subscriptionName,
	deadLetterSubscriptionName,
	host string
	pollingTime int
	opts        []nats.Option
}

// NewNATSConfig is a standard constructor
func NewNATSConfig(
	usr, pwd, subsName, deadLetter, h string,
	p uint16, pollingTime int, opts ...nats.Option,
) NATSConfig {
	return NATSConfig{
		user:                       usr,
		pass:                       pwd,
		subscriptionName:           subsName,
		deadLetterSubscriptionName: deadLetter,
		host:                       h,
		port:                       p,
		pollingTime:                pollingTime,
		opts:                       opts,
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

const (
	deadLetter = "dead.%s"
	// 1 Second
	defaultPollingTime = 1
)

// DefaultNATSConfig returns a standard configuration for a barebones NATS server
// for the passed subscription's name
func DefaultNATSConfig(subsName string) NATSConfig {
	return NATSConfig{
		user:                       "",
		pass:                       "",
		subscriptionName:           subsName,
		deadLetterSubscriptionName: fmt.Sprintf(deadLetter, subsName),
		host:                       "",
		pollingTime:                defaultPollingTime,
		port:                       0,
	}
}

// NATSBroker implementation of a Broker on top of NATS (https://nats.io/)
type NATSBroker struct {
	bus                    EventBus
	conn                   *nats.Conn
	conf                   NATSConfig
	messagesChan           chan *nats.Msg
	deadLetterMessagesChan chan *nats.Msg
}

// NewNATSBroker is a standard constructor
func NewNATSBroker(bus EventBus, conf NATSConfig) (*NATSBroker, error) {
	conn, err := nats.Connect(conf.URL(), conf.opts...)
	if err != nil {
		return nil, wrapError(ErrNATSServerConnection, err.Error())
	}
	messagesChan := make(chan *nats.Msg)
	_, err = conn.Subscribe(conf.subscriptionName, func(msg *nats.Msg) {
		messagesChan <- msg
	})
	if err != nil {
		return nil, wrapError(ErrNATSubscription, err.Error())
	}
	deadLetterMessagesChan := make(chan *nats.Msg)
	_, err = conn.Subscribe(conf.deadLetterSubscriptionName, func(msg *nats.Msg) {
		deadLetterMessagesChan <- msg
	})
	if err != nil {
		return nil, wrapError(ErrNATSubscription, err.Error())
	}
	return &NATSBroker{
		bus:                    bus,
		conn:                   conn,
		conf:                   conf,
		messagesChan:           messagesChan,
		deadLetterMessagesChan: deadLetterMessagesChan,
	}, nil
}

// CloseConnection closes the underlying NATS connection;
// meant for cleanUp purposes
func (n *NATSBroker) CloseConnection() {
	n.conn.Close()
}

// Process starts cosuming messages from a given subscription
func (n *NATSBroker) Process(ctx context.Context) <-chan error {
	errChan := make(chan error)
	errMsgHandler := func(err error, msg *nats.Msg) {
		errChan <- wrapError(ErrEventBus, err.Error())
		err = n.conn.Publish(fmt.Sprintf(deadLetter, msg.Subject), msg.Data)
		if err != nil {
			errChan <- wrapError(ErrDeadLetterPublish, err.Error())
		}
	}
	errHandler := func(err error) {
		errChan <- err
	}
	go func() {
		defer close(errChan)
		n.processMsgChan(ctx, n.messagesChan, errHandler, errMsgHandler)
	}()
	return errChan
}

// ProcessDead starts cosuming messages from a given subscription's deadLetter
func (n *NATSBroker) ProcessDead(ctx context.Context) <-chan error {
	errChan := make(chan error)
	errHandler := func(err error) {
		errChan <- err
	}
	errMsgHandler := func(err error, _ *nats.Msg) {
		errHandler(err)
	}
	go func() {
		defer close(errChan)
		n.processMsgChan(ctx, n.messagesChan, errHandler, errMsgHandler)
	}()
	return errChan
}

func (n *NATSBroker) processMsgChan(
	ctx context.Context, msgChan chan *nats.Msg,
	errHandler func(err error),
	errMsgHandler func(err error, msg *nats.Msg),
) {
	pollTicker := time.NewTicker(time.Duration(n.conf.pollingTime) * time.Second)
	for {
		defer pollTicker.Stop()
		select {
		case <-ctx.Done():
			errHandler(wrapError(ErrContextCanceled, ctx.Err().Error()))
			return
		case <-pollTicker.C:
			err := n.conn.Flush()
			if err != nil {
				errHandler(wrapError(ErrFlushSubscription, err.Error()))
				return
			}
		case msg := <-msgChan:
			event := Message{
				data: msg.Data,
			}
			err := n.bus.Resolve(ctx, event)
			if err != nil {
				errMsgHandler(err, msg)
			}
		}
	}
}

// Message is a wrapper for *nats.Msg
type Message struct {
	data []byte
}

// Data returns the data associated to the Message
func (m Message) Data() []byte { return m.data }

func wrapError(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

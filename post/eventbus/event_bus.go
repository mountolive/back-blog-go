package eventbus

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrEventNotRegistered = errors.New("event not registered")
	ErrCommandHandler     = errors.New("command handler error")
)

// Message describes a command's message
type Params interface{}

// Event describes a basic event in the application
type Event interface {
	Name() string
	Params() Params
}

// CommandHandler refers to capabilities that changes state in the application
type CommandHandler interface {
	Handle(context.Context, Params) error
}

// EventBus has a registry of Events against CommandHandlers
type EventBus struct {
	handlers map[string]CommandHandler
}

// NewEventBus creates an EventHandler
func NewEventBus() *EventBus {
	return &EventBus{handlers: make(map[string]CommandHandler)}
}

// Resolve passes an event using its corresponding CommandHandler
func (e EventBus) Resolve(ctx context.Context, event Event) error {
	handler, ok := e.handlers[event.Name()]
	if !ok {
		return ErrEventNotRegistered
	}
	err := handler.Handle(ctx, event)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCommandHandler, err)
	}
	return nil
}

// Register associates an Event with a given CommandHandler
func (e *EventBus) Register(eventName string, cmdHandler CommandHandler) {
	e.handlers[eventName] = cmdHandler
}

package eventbus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrEventNotRegistered is self-described
	ErrEventNotRegistered = errors.New("event not registered")
	// ErrEventNotRegistered is returned when an error is returned by a CommandHandler
	ErrCommandHandler = errors.New("command handler error")
	// ErrUnmarshalingMessage is self-described
	ErrUnmarshalingMessage = errors.New("unmarshaling message error")
	// ErrMissingNameParam returned when the key "name" is missing from an Event
	ErrMissingNameParam = errors.New("missing `name` param from message")
	// ErrWrongDataTypeName is self-described
	ErrWrongDataTypeName = errors.New("wrong data type for param `name`")
)

// Params is a map of the parameters of a call
type Params map[string]interface{}

// Event describes a basic event in the application
type Event interface {
	Data() []byte
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

// Resolve passes an event and executes its corresponding CommandHandler
func (e EventBus) Resolve(ctx context.Context, event Event) error {
	decodedEvent := make(map[string]interface{})
	eventData := strings.ReplaceAll(string(event.Data()), string('\x00'), "")
	err := json.Unmarshal([]byte(eventData), &decodedEvent)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUnmarshalingMessage, err)
	}
	nameResult, ok := decodedEvent["event_name"]
	if !ok {
		return ErrMissingNameParam
	}
	name, ok := nameResult.(string)
	if !ok {
		return ErrWrongDataTypeName
	}
	handler, ok := e.handlers[name]
	if !ok {
		return ErrEventNotRegistered
	}
	delete(decodedEvent, "event_name")
	err = handler.Handle(ctx, decodedEvent)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCommandHandler, err)
	}
	return nil
}

// Register associates an Event with a given CommandHandler
func (e *EventBus) Register(eventName string, cmdHandler CommandHandler) {
	e.handlers[eventName] = cmdHandler
}

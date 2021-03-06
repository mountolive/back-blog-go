package eventbus

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventHandler(t *testing.T) {
	t.Run("NewEventHandler", func(t *testing.T) {
		t.Parallel()
		bus := NewEventBus()
		require.NotNil(t, bus)
	})

	t.Run("Register", func(t *testing.T) {
		t.Parallel()
		bus := NewEventBus()
		eventName := "life on mars"
		var cmdHandler CommandHandler
		cmdHandler = &mockErroredCommandHandler{}
		bus.Register(eventName, cmdHandler)
		event := &testEvent{name: eventName}
		ctx := context.Background()
		err := bus.Resolve(ctx, event)
		require.True(t, errors.Is(err, ErrCommandHandler))
		cmdHandler = &mockCommandHandler{}
		bus.Register(eventName, cmdHandler)
		err = bus.Resolve(ctx, event)
		require.NoError(t, err)
	})

	t.Run("Resolve", func(t *testing.T) {
		t.Parallel()
		bus := NewEventBus()
		err := bus.Resolve(context.Background(), &testEvent{name: "not-existing"})
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrEventNotRegistered))
	})
}

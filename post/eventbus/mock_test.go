package eventbus

import (
	"context"
	"fmt"
)

var _ CommandHandler = &mockCommandHandler{}

type mockCommandHandler struct{}

func (*mockCommandHandler) Handle(_ context.Context, p Params) error {
	return nil
}

var _ CommandHandler = &mockErroredCommandHandler{}

type mockErroredCommandHandler struct{}

func (*mockErroredCommandHandler) Handle(_ context.Context, p Params) error {
	return fmt.Errorf("an error occurred in CommandHandler")
}

var _ Event = &testEvent{}

type testEvent struct {
	name string
}

func (e testEvent) Data() []byte {
	return []byte(
		fmt.Sprintf(`{"event_name": "%s", "data": "some-data"}`, e.name),
	)
}

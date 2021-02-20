package handlers

import "testing"

func TestCreatePostCommandHandler(t *testing.T) {
	t.Run("Canary", func(t *testing.T) {
		var _ CommandHandler = CreatePostCommandHandler{}
	})

	t.Run("Handle", func(t *testing.T) {
	})
}

func TestUpdatePostCommandHandler(t *testing.T) {
	t.Run("Canary", func(t *testing.T) {
		var _ CommandHandler = UpdatePostCommandHandler{}
	})

	t.Run("Handle", func(t *testing.T) {
	})
}

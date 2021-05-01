package command_test

import (
	"context"

	"github.com/mountolive/back-blog-go/post/usecase"
)

type mockStore struct{}

func (*mockStore) Create(context.Context, *usecase.CreatePostDto) (*usecase.Post, error) {
	return nil, nil
}

func (*mockStore) Update(context.Context, *usecase.UpdatePostDto) (*usecase.Post, error) {
	return nil, nil
}

func (*mockStore) Filter(context.Context, *usecase.GeneralFilter) ([]*usecase.Post, error) {
	return nil, nil
}

func (*mockStore) ReadOne(context.Context, string) (*usecase.Post, error) {
	return nil, nil
}

type mockStoreErrored struct {
	err error
}

func (m *mockStoreErrored) Create(context.Context, *usecase.CreatePostDto) (*usecase.Post, error) {
	return nil, m.err
}

func (m *mockStoreErrored) Update(context.Context, *usecase.UpdatePostDto) (*usecase.Post, error) {
	return nil, m.err
}

func (*mockStoreErrored) Filter(context.Context, *usecase.GeneralFilter) ([]*usecase.Post, error) {
	return nil, nil
}

func (*mockStoreErrored) ReadOne(context.Context, string) (*usecase.Post, error) {
	return nil, nil
}

type mockTrueChecker struct{}

func (m *mockTrueChecker) CheckExistence(ctx context.Context, c string) (bool, error) {
	return true, nil
}

type mockFalseChecker struct{}

func (m *mockFalseChecker) CheckExistence(ctx context.Context, c string) (bool, error) {
	return false, nil
}

type mockErroredChecker struct {
	err error
}

func (m *mockErroredChecker) CheckExistence(ctx context.Context, c string) (bool, error) {
	return false, m.err
}

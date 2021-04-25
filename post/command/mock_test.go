package command_test

import (
	"context"
)

type mockStore struct{}

func (*mockStore) Create(context.Context, *CreatePostDto) (*Post, error) {
	return nil, nil
}

func (*mockStore) Update(context.Context, *UpdatePostDto) (*Post, error) {
	return nil, nil
}

func (*mockStore) Filter(context.Context, *GeneralFilter) ([]*Post, error) {
	return nil, nil
}

func (*mockStore) ReadOne(context.Context, string) (*Post, error) {
	return nil, nil
}

type mockStoreErrored struct {
	err error
}

func (m *mockStoreErrored) Create(context.Context, *CreatePostDto) (*Post, error) {
	return nil, m.err
}

func (m *mockStoreErrored) Update(context.Context, *UpdatePostDto) (*Post, error) {
	return nil, m.err
}

func (*mockStoreErrored) Filter(context.Context, *GeneralFilter) ([]*Post, error) {
	return nil, nil
}

func (*mockStoreErrored) ReadOne(context.Context, string) (*Post, error) {
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
	return false, err
}

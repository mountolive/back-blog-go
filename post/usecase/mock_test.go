package usecase

import (
	"context"
	"errors"
	"fmt"
)

// Mock logger
type mockLogger struct{}

func (m *mockLogger) LogError(err error) {
	fmt.Printf("%v \n", err)
}

type mockStoreNotEmpty struct{}

func (m *mockStoreNotEmpty) Create(ctx context.Context, p *CreatePostDto) (*Post, error) {
	return &Post{Creator: p.Creator, Content: p.Content}, nil
}

func (m *mockStoreNotEmpty) Update(ctx context.Context, p *UpdatePostDto) (*Post, error) {
	return &Post{Id: p.Id, Creator: "test", Content: p.Content}, nil
}

func (m *mockStoreNotEmpty) Filter(ctx context.Context, p *GeneralFilter) ([]*Post, error) {
	return []*Post{&Post{Creator: "test", Content: "test", Tags: []string{p.Tag}}}, nil
}

func (m *mockStoreNotEmpty) ReadOne(ctx context.Context, id string) (*Post, error) {
	return &Post{Id: id, Creator: "bla", Content: "hello"}, nil
}

type mockStoreEmpty struct{}

func (m *mockStoreEmpty) Create(ctx context.Context, p *CreatePostDto) (*Post, error) {
	return nil, nil
}

func (m *mockStoreEmpty) Update(ctx context.Context, p *UpdatePostDto) (*Post, error) {
	return nil, errors.New("Any error occurred")
}

func (m *mockStoreEmpty) Filter(ctx context.Context, p *GeneralFilter) ([]*Post, error) {
	return nil, nil
}

func (m *mockStoreEmpty) ReadOne(ctx context.Context, id string) (*Post, error) {
	return &Post{}, nil
}

type mockStoreReadErrored struct{}

func (m *mockStoreReadErrored) Create(ctx context.Context, p *CreatePostDto) (*Post, error) {
	return nil, nil
}

func (m *mockStoreReadErrored) Update(ctx context.Context, p *UpdatePostDto) (*Post, error) {
	return nil, errors.New("Any error occurred")
}

func (m *mockStoreReadErrored) Filter(ctx context.Context, p *GeneralFilter) ([]*Post, error) {
	return nil, nil
}

func (m *mockStoreReadErrored) ReadOne(ctx context.Context, id string) (*Post, error) {
	return nil, errors.New("Something happened")
}

type mockSanitizer struct{}

func (m *mockSanitizer) SanitizeContent(content string) string {
	return content
}

type mockFalseChecker struct{}

func (m *mockFalseChecker) CheckExistence(ctx context.Context, c string) (bool, error) {
	return false, nil
}

type mockTrueChecker struct{}

func (m *mockTrueChecker) CheckExistence(ctx context.Context, c string) (bool, error) {
	return true, nil
}

type mockErrorChecker struct{}

func (m *mockErrorChecker) CheckExistence(ctx context.Context, c string) (bool, error) {
	return false, errors.New("Not found")
}

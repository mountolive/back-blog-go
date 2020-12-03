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

func (m *mockStoreNotEmpty) Create(ctx context.Context, p *CreatePostDto) (*PostDto, error) {
	return &PostDto{Creator: p.Creator, Content: p.Content}, nil
}

func (m *mockStoreNotEmpty) Update(ctx context.Context, p *UpdatePostDto) (*PostDto, error) {
	return &PostDto{Id: p.Id, Creator: "test", Content: p.Content}, nil
}

func (m *mockStoreNotEmpty) Filter(ctx context.Context, p *GeneralFilter) ([]*PostDto, error) {
	return []*PostDto{&PostDto{Creator: "test", Content: "test", Tags: []Tag{Tag{p.Tag}}}}, nil
}

func (m *mockStoreNotEmpty) ReadOne(ctx context.Context, id string) (*PostDto, error) {
	return &PostDto{Id: id, Creator: "bla", Content: "hello"}, nil
}

type mockStoreEmpty struct{}

func (m *mockStoreEmpty) Create(ctx context.Context, p *CreatePostDto) (*PostDto, error) {
	return nil, nil
}

func (m *mockStoreEmpty) Update(ctx context.Context, p *UpdatePostDto) (*PostDto, error) {
	return nil, errors.New("Any error occurred")
}

func (m *mockStoreEmpty) Filter(ctx context.Context, p *GeneralFilter) ([]*PostDto, error) {
	return nil, nil
}

func (m *mockStoreEmpty) ReadOne(ctx context.Context, id string) (*PostDto, error) {
	return nil, errors.New("Any error occurred")
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

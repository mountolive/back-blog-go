package user

import (
	"context"
	"fmt"

	"github.com/mountolive/back-blog-go/post/usecase"
)

// TODO Add test suite for GRPCUserChecker, posts service

// GRPCUserChecker wraps a gRPC client to connect to users service
type GRPCUserChecker struct {
	client UserClient
}

// NewGRPCUserChecker is a constructor
func NewGRPCUserChecker(client UserClient) GRPCUserChecker {
	return GRPCUserChecker{client}
}

var _ usecase.CreatorChecker = GRPCUserChecker{}

const errMsgCheckUser = "grpc client check user: %w"

// CheckExistence implements CreatorChecker interface
func (g GRPCUserChecker) CheckExistence(ctx context.Context, login string) (bool, error) {
	req := &CheckUserRequest{Login: login}
	res, err := g.client.CheckUser(ctx, req)
	if err != nil {
		return false, fmt.Errorf(errMsgCheckUser, err)
	}
	if res.Id == "" {
		return false, nil
	}
	return true, nil
}

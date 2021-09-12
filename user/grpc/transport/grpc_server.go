package transport

import (
	"context"
	"fmt"

	"github.com/mountolive/back-blog-go/user/usecase"
)

// TODO Add test suite for GRPCServer, users service

const (
	errMsgCreateUser     = "grpc create user: %w"
	errMsgUpdateUser     = "grpc update user: %w"
	errMsgChangePassword = "grpc change password from user: %w"
	errMsgCheckUser      = "grpc check user: %w"
)

// GRPCServer is self-described
type GRPCServer struct {
	UnimplementedUserCheckerServer
	UnimplementedUserCreatorServer
	UnimplementedUserUpdaterServer
	UnimplementedPasswordChangerServer
	UnimplementedLoginServer
	repo usecase.Repository
}

// NewGRPCServer is a constructor
func NewGRPCServer(repo usecase.Repository) GRPCServer {
	return GRPCServer{
		UnimplementedUserCheckerServer{},
		UnimplementedUserCreatorServer{},
		UnimplementedUserUpdaterServer{},
		UnimplementedPasswordChangerServer{},
		UnimplementedLoginServer{},
		repo,
	}
}

var (
	_ UserCheckerServer     = GRPCServer{}
	_ UserCreatorServer     = GRPCServer{}
	_ UserUpdaterServer     = GRPCServer{}
	_ PasswordChangerServer = GRPCServer{}
	_ LoginServer           = GRPCServer{}
)

func newCreateUserDto(cu *CreateUserRequest) *usecase.CreateUserDto {
	return &usecase.CreateUserDto{
		Email:            cu.Email,
		Username:         cu.Username,
		Password:         cu.Password,
		RepeatedPassword: cu.RepeatedPassword,
		FirstName:        cu.FirstName,
		LastName:         cu.LastName,
	}
}

// Create implements the UserServer interface
func (g GRPCServer) Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	user, err := g.repo.CreateUser(ctx, newCreateUserDto(req))
	if err != nil {
		return nil, fmt.Errorf(errMsgCreateUser, err)
	}
	return newUserResponse(user), nil
}

func newUpdateUserDto(uu *UpdateUserRequest) *usecase.UpdateUserDto {
	return &usecase.UpdateUserDto{
		Email:     uu.Email,
		Username:  uu.Username,
		FirstName: uu.FirstName,
		LastName:  uu.LastName,
	}
}

// Update implements the UserUpdaterServer interface
func (g GRPCServer) Update(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	user, err := g.repo.UpdateUser(ctx, req.Id, newUpdateUserDto(req))
	if err != nil {
		return nil, fmt.Errorf(errMsgUpdateUser, err)
	}
	return newUserResponse(user), nil
}

func newChangePasswordDto(cp *ChangePasswordRequest) *usecase.ChangePasswordDto {
	return &usecase.ChangePasswordDto{
		Email:            cp.Email,
		Username:         cp.Username,
		OldPassword:      cp.OldPassword,
		NewPassword:      cp.NewPassword,
		RepeatedPassword: cp.RepeatedPassword,
	}
}

// ChangePassword implements the PasswordChangerServer interface
func (g GRPCServer) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	err := g.repo.ChangePassword(ctx, newChangePasswordDto(req))
	if err != nil {
		return nil, fmt.Errorf(errMsgChangePassword, err)
	}
	return &ChangePasswordResponse{Success: true}, nil
}

// CheckUser implements the UserCheckerServer interface
func (g GRPCServer) CheckUser(ctx context.Context, req *CheckUserRequest) (*UserResponse, error) {
	user, err := g.repo.ReadUser(ctx, req.Login)
	if err != nil {
		return nil, fmt.Errorf(errMsgCheckUser, err)
	}
	return newUserResponse(user), nil
}

// Login implements the LoginServer interface
func (g GRPCServer) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// TODO: Implement Login method in users' repo and then in GRPCServer
	return nil, nil
}

func newUserResponse(u *usecase.User) *UserResponse {
	return &UserResponse{
		Id:        u.Id,
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

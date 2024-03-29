// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package transport

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserCreatorClient is the client API for UserCreator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserCreatorClient interface {
	Create(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
}

type userCreatorClient struct {
	cc grpc.ClientConnInterface
}

func NewUserCreatorClient(cc grpc.ClientConnInterface) UserCreatorClient {
	return &userCreatorClient{cc}
}

func (c *userCreatorClient) Create(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/user.UserCreator/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserCreatorServer is the server API for UserCreator service.
// All implementations must embed UnimplementedUserCreatorServer
// for forward compatibility
type UserCreatorServer interface {
	Create(context.Context, *CreateUserRequest) (*UserResponse, error)
	mustEmbedUnimplementedUserCreatorServer()
}

// UnimplementedUserCreatorServer must be embedded to have forward compatible implementations.
type UnimplementedUserCreatorServer struct {
}

func (UnimplementedUserCreatorServer) Create(context.Context, *CreateUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedUserCreatorServer) mustEmbedUnimplementedUserCreatorServer() {}

// UnsafeUserCreatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserCreatorServer will
// result in compilation errors.
type UnsafeUserCreatorServer interface {
	mustEmbedUnimplementedUserCreatorServer()
}

func RegisterUserCreatorServer(s grpc.ServiceRegistrar, srv UserCreatorServer) {
	s.RegisterService(&UserCreator_ServiceDesc, srv)
}

func _UserCreator_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserCreatorServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserCreator/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserCreatorServer).Create(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserCreator_ServiceDesc is the grpc.ServiceDesc for UserCreator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserCreator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserCreator",
	HandlerType: (*UserCreatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _UserCreator_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "creator.proto",
}

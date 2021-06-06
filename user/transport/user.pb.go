// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transport/user.proto

package transport

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type UserResponse struct {
	Id                   string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Email                string               `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Username             string               `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	FirstName            string               `protobuf:"bytes,4,opt,name=firstName,proto3" json:"firstName,omitempty"`
	LastName             string               `protobuf:"bytes,5,opt,name=lastName,proto3" json:"lastName,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,6,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,7,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *UserResponse) Reset()         { *m = UserResponse{} }
func (m *UserResponse) String() string { return proto.CompactTextString(m) }
func (*UserResponse) ProtoMessage()    {}
func (*UserResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6fa1ee495e6d9f92, []int{0}
}

func (m *UserResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserResponse.Unmarshal(m, b)
}
func (m *UserResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserResponse.Marshal(b, m, deterministic)
}
func (m *UserResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserResponse.Merge(m, src)
}
func (m *UserResponse) XXX_Size() int {
	return xxx_messageInfo_UserResponse.Size(m)
}
func (m *UserResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UserResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UserResponse proto.InternalMessageInfo

func (m *UserResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UserResponse) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UserResponse) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *UserResponse) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *UserResponse) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *UserResponse) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *UserResponse) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

// Empty response for ChangePassword
type ChangePasswordResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangePasswordResponse) Reset()         { *m = ChangePasswordResponse{} }
func (m *ChangePasswordResponse) String() string { return proto.CompactTextString(m) }
func (*ChangePasswordResponse) ProtoMessage()    {}
func (*ChangePasswordResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_6fa1ee495e6d9f92, []int{1}
}

func (m *ChangePasswordResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChangePasswordResponse.Unmarshal(m, b)
}
func (m *ChangePasswordResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChangePasswordResponse.Marshal(b, m, deterministic)
}
func (m *ChangePasswordResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangePasswordResponse.Merge(m, src)
}
func (m *ChangePasswordResponse) XXX_Size() int {
	return xxx_messageInfo_ChangePasswordResponse.Size(m)
}
func (m *ChangePasswordResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangePasswordResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ChangePasswordResponse proto.InternalMessageInfo

func (m *ChangePasswordResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type CreateUserRequest struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Username             string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	RepeatedPassword     string   `protobuf:"bytes,4,opt,name=repeatedPassword,proto3" json:"repeatedPassword,omitempty"`
	FirstName            string   `protobuf:"bytes,5,opt,name=firstName,proto3" json:"firstName,omitempty"`
	LastName             string   `protobuf:"bytes,6,opt,name=lastName,proto3" json:"lastName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateUserRequest) Reset()         { *m = CreateUserRequest{} }
func (m *CreateUserRequest) String() string { return proto.CompactTextString(m) }
func (*CreateUserRequest) ProtoMessage()    {}
func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6fa1ee495e6d9f92, []int{2}
}

func (m *CreateUserRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateUserRequest.Unmarshal(m, b)
}
func (m *CreateUserRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateUserRequest.Marshal(b, m, deterministic)
}
func (m *CreateUserRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateUserRequest.Merge(m, src)
}
func (m *CreateUserRequest) XXX_Size() int {
	return xxx_messageInfo_CreateUserRequest.Size(m)
}
func (m *CreateUserRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateUserRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateUserRequest proto.InternalMessageInfo

func (m *CreateUserRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *CreateUserRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *CreateUserRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *CreateUserRequest) GetRepeatedPassword() string {
	if m != nil {
		return m.RepeatedPassword
	}
	return ""
}

func (m *CreateUserRequest) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *CreateUserRequest) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

type UpdateUserRequest struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Username             string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	FirstName            string   `protobuf:"bytes,3,opt,name=firstName,proto3" json:"firstName,omitempty"`
	LastName             string   `protobuf:"bytes,4,opt,name=lastName,proto3" json:"lastName,omitempty"`
	Id                   string   `protobuf:"bytes,5,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateUserRequest) Reset()         { *m = UpdateUserRequest{} }
func (m *UpdateUserRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateUserRequest) ProtoMessage()    {}
func (*UpdateUserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6fa1ee495e6d9f92, []int{3}
}

func (m *UpdateUserRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateUserRequest.Unmarshal(m, b)
}
func (m *UpdateUserRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateUserRequest.Marshal(b, m, deterministic)
}
func (m *UpdateUserRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateUserRequest.Merge(m, src)
}
func (m *UpdateUserRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateUserRequest.Size(m)
}
func (m *UpdateUserRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateUserRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateUserRequest proto.InternalMessageInfo

func (m *UpdateUserRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *UpdateUserRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *UpdateUserRequest) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *UpdateUserRequest) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *UpdateUserRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type ChangePasswordRequest struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Username             string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	NewPassword          string   `protobuf:"bytes,3,opt,name=newPassword,proto3" json:"newPassword,omitempty"`
	RepeatedPassword     string   `protobuf:"bytes,4,opt,name=repeatedPassword,proto3" json:"repeatedPassword,omitempty"`
	OldPassword          string   `protobuf:"bytes,5,opt,name=oldPassword,proto3" json:"oldPassword,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangePasswordRequest) Reset()         { *m = ChangePasswordRequest{} }
func (m *ChangePasswordRequest) String() string { return proto.CompactTextString(m) }
func (*ChangePasswordRequest) ProtoMessage()    {}
func (*ChangePasswordRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6fa1ee495e6d9f92, []int{4}
}

func (m *ChangePasswordRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChangePasswordRequest.Unmarshal(m, b)
}
func (m *ChangePasswordRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChangePasswordRequest.Marshal(b, m, deterministic)
}
func (m *ChangePasswordRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangePasswordRequest.Merge(m, src)
}
func (m *ChangePasswordRequest) XXX_Size() int {
	return xxx_messageInfo_ChangePasswordRequest.Size(m)
}
func (m *ChangePasswordRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangePasswordRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ChangePasswordRequest proto.InternalMessageInfo

func (m *ChangePasswordRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *ChangePasswordRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *ChangePasswordRequest) GetNewPassword() string {
	if m != nil {
		return m.NewPassword
	}
	return ""
}

func (m *ChangePasswordRequest) GetRepeatedPassword() string {
	if m != nil {
		return m.RepeatedPassword
	}
	return ""
}

func (m *ChangePasswordRequest) GetOldPassword() string {
	if m != nil {
		return m.OldPassword
	}
	return ""
}

type CheckUserRequest struct {
	Login                string   `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckUserRequest) Reset()         { *m = CheckUserRequest{} }
func (m *CheckUserRequest) String() string { return proto.CompactTextString(m) }
func (*CheckUserRequest) ProtoMessage()    {}
func (*CheckUserRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6fa1ee495e6d9f92, []int{5}
}

func (m *CheckUserRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckUserRequest.Unmarshal(m, b)
}
func (m *CheckUserRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckUserRequest.Marshal(b, m, deterministic)
}
func (m *CheckUserRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckUserRequest.Merge(m, src)
}
func (m *CheckUserRequest) XXX_Size() int {
	return xxx_messageInfo_CheckUserRequest.Size(m)
}
func (m *CheckUserRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckUserRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CheckUserRequest proto.InternalMessageInfo

func (m *CheckUserRequest) GetLogin() string {
	if m != nil {
		return m.Login
	}
	return ""
}

func init() {
	proto.RegisterType((*UserResponse)(nil), "transport.UserResponse")
	proto.RegisterType((*ChangePasswordResponse)(nil), "transport.ChangePasswordResponse")
	proto.RegisterType((*CreateUserRequest)(nil), "transport.CreateUserRequest")
	proto.RegisterType((*UpdateUserRequest)(nil), "transport.UpdateUserRequest")
	proto.RegisterType((*ChangePasswordRequest)(nil), "transport.ChangePasswordRequest")
	proto.RegisterType((*CheckUserRequest)(nil), "transport.CheckUserRequest")
}

func init() {
	proto.RegisterFile("transport/user.proto", fileDescriptor_6fa1ee495e6d9f92)
}

var fileDescriptor_6fa1ee495e6d9f92 = []byte{
	// 458 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0xc7, 0x71, 0xd6, 0x76, 0xcd, 0x29, 0x4c, 0x9b, 0x35, 0x20, 0x0a, 0x93, 0x08, 0xb9, 0xaa,
	0xb8, 0x48, 0xa5, 0x72, 0xc3, 0xed, 0xe8, 0x3d, 0x9a, 0x2a, 0x26, 0x24, 0xee, 0xbc, 0xe6, 0xac,
	0x8b, 0x48, 0xe3, 0x60, 0x3b, 0xda, 0x63, 0xf0, 0x1e, 0xbc, 0x01, 0x4f, 0xc0, 0x5b, 0x21, 0x14,
	0xbb, 0x76, 0xe3, 0x16, 0x82, 0x60, 0x97, 0xe7, 0xd3, 0xff, 0xf3, 0x3b, 0xc7, 0x70, 0xae, 0x04,
	0xab, 0x64, 0xcd, 0x85, 0x9a, 0x35, 0x12, 0x45, 0x56, 0x0b, 0xae, 0x38, 0x0d, 0x9d, 0x37, 0x7e,
	0xb9, 0xe6, 0x7c, 0x5d, 0xe2, 0x4c, 0x07, 0x6e, 0x9a, 0xdb, 0x99, 0x2a, 0x36, 0x28, 0x15, 0xdb,
	0xd4, 0x26, 0x37, 0xfd, 0x49, 0xe0, 0xf1, 0xb5, 0x44, 0xb1, 0x44, 0x59, 0xf3, 0x4a, 0x22, 0x3d,
	0x81, 0xa0, 0xc8, 0x23, 0x92, 0x90, 0x69, 0xb8, 0x0c, 0x8a, 0x9c, 0x9e, 0xc3, 0x10, 0x37, 0xac,
	0x28, 0xa3, 0x40, 0xbb, 0x8c, 0x41, 0x63, 0x18, 0xb7, 0x0f, 0x56, 0x6c, 0x83, 0xd1, 0x91, 0x0e,
	0x38, 0x9b, 0x5e, 0x40, 0x78, 0x5b, 0x08, 0xa9, 0xde, 0xb7, 0xc1, 0x81, 0x0e, 0xee, 0x1c, 0x6d,
	0x65, 0xc9, 0xb6, 0xc1, 0xa1, 0xa9, 0xb4, 0x36, 0x7d, 0x0b, 0xe1, 0x4a, 0x20, 0x53, 0x98, 0x5f,
	0xaa, 0x68, 0x94, 0x90, 0xe9, 0x64, 0x1e, 0x67, 0x66, 0x82, 0xcc, 0x4e, 0x90, 0x7d, 0xb0, 0x13,
	0x2c, 0x77, 0xc9, 0x6d, 0x65, 0x53, 0xe7, 0xdb, 0xca, 0xe3, 0xbf, 0x57, 0xba, 0xe4, 0x74, 0x0e,
	0xcf, 0x16, 0x77, 0xac, 0x5a, 0xe3, 0x15, 0x93, 0xf2, 0x9e, 0x8b, 0xdc, 0x91, 0x88, 0xe0, 0x58,
	0x36, 0xab, 0x15, 0x4a, 0xa9, 0x71, 0x8c, 0x97, 0xd6, 0x4c, 0x7f, 0x10, 0x38, 0x5b, 0xe8, 0xb7,
	0x0d, 0xba, 0x2f, 0x0d, 0x4a, 0xb5, 0x23, 0x45, 0xfe, 0x44, 0x2a, 0xd8, 0x23, 0x15, 0xc3, 0xb8,
	0xde, 0xbe, 0x6a, 0x29, 0x5a, 0x9b, 0xbe, 0x86, 0x53, 0x81, 0xb5, 0x9e, 0xcf, 0x2a, 0xdb, 0xc2,
	0x3c, 0xf0, 0xfb, 0xc4, 0x87, 0x7d, 0xc4, 0x47, 0x3e, 0xf1, 0xf4, 0x2b, 0x81, 0xb3, 0x6b, 0xcd,
	0xe2, 0x61, 0x93, 0x78, 0x0a, 0x8e, 0xfa, 0x14, 0x0c, 0xf6, 0x76, 0x6e, 0xee, 0x6d, 0x68, 0xef,
	0x2d, 0xfd, 0x4e, 0xe0, 0xe9, 0xfe, 0x42, 0xfe, 0x57, 0x55, 0x02, 0x93, 0x0a, 0xef, 0xaf, 0x7c,
	0xc4, 0x5d, 0xd7, 0x3f, 0x51, 0x4e, 0x60, 0xc2, 0xcb, 0x5d, 0x9a, 0x91, 0xdc, 0x75, 0xa5, 0x53,
	0x38, 0x5d, 0xdc, 0xe1, 0xea, 0xf3, 0x1e, 0xcb, 0x92, 0xaf, 0x8b, 0xca, 0xaa, 0xd6, 0xc6, 0xfc,
	0x5b, 0x00, 0x83, 0x36, 0x8b, 0x5e, 0xc2, 0xc8, 0x5c, 0x12, 0xbd, 0xc8, 0xdc, 0xb7, 0xcd, 0x0e,
	0x8e, 0x2b, 0x7e, 0xde, 0x89, 0x76, 0xff, 0x6b, 0xfa, 0xa8, 0x6d, 0x61, 0x56, 0xe8, 0xb5, 0x38,
	0xd8, 0x6a, 0x5f, 0x8b, 0x8f, 0x70, 0xe2, 0x33, 0xa7, 0x49, 0x57, 0xcd, 0xef, 0xd6, 0x11, 0xbf,
	0xea, 0xc9, 0x70, 0x8d, 0x17, 0x10, 0x3a, 0x22, 0xf4, 0x85, 0x57, 0xe1, 0x73, 0xea, 0x51, 0xf7,
	0xee, 0xc9, 0xa7, 0x49, 0x36, 0x73, 0xd1, 0x9b, 0x91, 0xfe, 0xd0, 0x6f, 0x7e, 0x05, 0x00, 0x00,
	0xff, 0xff, 0xe7, 0x30, 0xa8, 0x11, 0xfd, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserClient interface {
	Create(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	Update(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
	ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error)
	CheckUser(ctx context.Context, in *CheckUserRequest, opts ...grpc.CallOption) (*UserResponse, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) Create(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/transport.User/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Update(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/transport.User/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error) {
	out := new(ChangePasswordResponse)
	err := c.cc.Invoke(ctx, "/transport.User/ChangePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) CheckUser(ctx context.Context, in *CheckUserRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/transport.User/CheckUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
type UserServer interface {
	Create(context.Context, *CreateUserRequest) (*UserResponse, error)
	Update(context.Context, *UpdateUserRequest) (*UserResponse, error)
	ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error)
	CheckUser(context.Context, *CheckUserRequest) (*UserResponse, error)
}

// UnimplementedUserServer can be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (*UnimplementedUserServer) Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedUserServer) Update(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (*UnimplementedUserServer) ChangePassword(ctx context.Context, req *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (*UnimplementedUserServer) CheckUser(ctx context.Context, req *CheckUserRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckUser not implemented")
}

func RegisterUserServer(s *grpc.Server, srv UserServer) {
	s.RegisterService(&_User_serviceDesc, srv)
}

func _User_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transport.User/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Create(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transport.User/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Update(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transport.User/ChangePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).ChangePassword(ctx, req.(*ChangePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_CheckUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).CheckUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transport.User/CheckUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).CheckUser(ctx, req.(*CheckUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _User_serviceDesc = grpc.ServiceDesc{
	ServiceName: "transport.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _User_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _User_Update_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _User_ChangePassword_Handler,
		},
		{
			MethodName: "CheckUser",
			Handler:    _User_CheckUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "transport/user.proto",
}

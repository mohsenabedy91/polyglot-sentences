// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: internal/adapter/grpc/proto/user/user.proto

package user

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	// Retrieves user information by UUID.
	GetByUUID(ctx context.Context, in *GetByUUIDRequest, opts ...grpc.CallOption) (*UserResponse, error)
	// Retrieves user information by email.
	GetByEmail(ctx context.Context, in *GetByEmailRequest, opts ...grpc.CallOption) (*UserResponse, error)
	// Checks if an email is unique.
	IsEmailUnique(ctx context.Context, in *IsEmailUniqueRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Creates a new user.
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*UserResponse, error)
	// Marks the user's email as verified.
	VerifiedEmail(ctx context.Context, in *VerifiedEmailRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Updates the flag indicating whether the welcome message was sent.
	MarkWelcomeMessageSent(ctx context.Context, in *UpdateWelcomeMessageToSentRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Updates the user google ID.
	UpdateGoogleID(ctx context.Context, in *UpdateGoogleIDRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	// Updates the user lat login time.
	UpdateLastLoginTime(ctx context.Context, in *UpdateLastLoginTimeRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) GetByUUID(ctx context.Context, in *GetByUUIDRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/user.UserService/GetByUUID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetByEmail(ctx context.Context, in *GetByEmailRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/user.UserService/GetByEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) IsEmailUnique(ctx context.Context, in *IsEmailUniqueRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/user.UserService/IsEmailUnique", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*UserResponse, error) {
	out := new(UserResponse)
	err := c.cc.Invoke(ctx, "/user.UserService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) VerifiedEmail(ctx context.Context, in *VerifiedEmailRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/user.UserService/VerifiedEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) MarkWelcomeMessageSent(ctx context.Context, in *UpdateWelcomeMessageToSentRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/user.UserService/MarkWelcomeMessageSent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateGoogleID(ctx context.Context, in *UpdateGoogleIDRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/user.UserService/UpdateGoogleID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UpdateLastLoginTime(ctx context.Context, in *UpdateLastLoginTimeRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/user.UserService/UpdateLastLoginTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	// Retrieves user information by UUID.
	GetByUUID(context.Context, *GetByUUIDRequest) (*UserResponse, error)
	// Retrieves user information by email.
	GetByEmail(context.Context, *GetByEmailRequest) (*UserResponse, error)
	// Checks if an email is unique.
	IsEmailUnique(context.Context, *IsEmailUniqueRequest) (*empty.Empty, error)
	// Creates a new user.
	Create(context.Context, *CreateRequest) (*UserResponse, error)
	// Marks the user's email as verified.
	VerifiedEmail(context.Context, *VerifiedEmailRequest) (*empty.Empty, error)
	// Updates the flag indicating whether the welcome message was sent.
	MarkWelcomeMessageSent(context.Context, *UpdateWelcomeMessageToSentRequest) (*empty.Empty, error)
	// Updates the user google ID.
	UpdateGoogleID(context.Context, *UpdateGoogleIDRequest) (*empty.Empty, error)
	// Updates the user lat login time.
	UpdateLastLoginTime(context.Context, *UpdateLastLoginTimeRequest) (*empty.Empty, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) GetByUUID(context.Context, *GetByUUIDRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByUUID not implemented")
}
func (UnimplementedUserServiceServer) GetByEmail(context.Context, *GetByEmailRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByEmail not implemented")
}
func (UnimplementedUserServiceServer) IsEmailUnique(context.Context, *IsEmailUniqueRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsEmailUnique not implemented")
}
func (UnimplementedUserServiceServer) Create(context.Context, *CreateRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedUserServiceServer) VerifiedEmail(context.Context, *VerifiedEmailRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifiedEmail not implemented")
}
func (UnimplementedUserServiceServer) MarkWelcomeMessageSent(context.Context, *UpdateWelcomeMessageToSentRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkWelcomeMessageSent not implemented")
}
func (UnimplementedUserServiceServer) UpdateGoogleID(context.Context, *UpdateGoogleIDRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGoogleID not implemented")
}
func (UnimplementedUserServiceServer) UpdateLastLoginTime(context.Context, *UpdateLastLoginTimeRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLastLoginTime not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_GetByUUID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByUUIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetByUUID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/GetByUUID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetByUUID(ctx, req.(*GetByUUIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/GetByEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetByEmail(ctx, req.(*GetByEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_IsEmailUnique_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsEmailUniqueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).IsEmailUnique(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/IsEmailUnique",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).IsEmailUnique(ctx, req.(*IsEmailUniqueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_VerifiedEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifiedEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).VerifiedEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/VerifiedEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).VerifiedEmail(ctx, req.(*VerifiedEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_MarkWelcomeMessageSent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateWelcomeMessageToSentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).MarkWelcomeMessageSent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/MarkWelcomeMessageSent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).MarkWelcomeMessageSent(ctx, req.(*UpdateWelcomeMessageToSentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateGoogleID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGoogleIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateGoogleID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/UpdateGoogleID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateGoogleID(ctx, req.(*UpdateGoogleIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UpdateLastLoginTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLastLoginTimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateLastLoginTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.UserService/UpdateLastLoginTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateLastLoginTime(ctx, req.(*UpdateLastLoginTimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByUUID",
			Handler:    _UserService_GetByUUID_Handler,
		},
		{
			MethodName: "GetByEmail",
			Handler:    _UserService_GetByEmail_Handler,
		},
		{
			MethodName: "IsEmailUnique",
			Handler:    _UserService_IsEmailUnique_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _UserService_Create_Handler,
		},
		{
			MethodName: "VerifiedEmail",
			Handler:    _UserService_VerifiedEmail_Handler,
		},
		{
			MethodName: "MarkWelcomeMessageSent",
			Handler:    _UserService_MarkWelcomeMessageSent_Handler,
		},
		{
			MethodName: "UpdateGoogleID",
			Handler:    _UserService_UpdateGoogleID_Handler,
		},
		{
			MethodName: "UpdateLastLoginTime",
			Handler:    _UserService_UpdateLastLoginTime_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/adapter/grpc/proto/user/user.proto",
}

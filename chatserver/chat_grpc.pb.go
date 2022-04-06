// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package chatserver

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ServicesClient is the client API for Services service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServicesClient interface {
	ChatService(ctx context.Context, opts ...grpc.CallOption) (Services_ChatServiceClient, error)
	CommandService(ctx context.Context, in *Command, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetClients(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Clients, error)
	VerifyName(ctx context.Context, in *ClientName, opts ...grpc.CallOption) (*ClientNameResponse, error)
}

type servicesClient struct {
	cc grpc.ClientConnInterface
}

func NewServicesClient(cc grpc.ClientConnInterface) ServicesClient {
	return &servicesClient{cc}
}

func (c *servicesClient) ChatService(ctx context.Context, opts ...grpc.CallOption) (Services_ChatServiceClient, error) {
	stream, err := c.cc.NewStream(ctx, &Services_ServiceDesc.Streams[0], "/pixalquarks.terminalChatServer.Services/ChatService", opts...)
	if err != nil {
		return nil, err
	}
	x := &servicesChatServiceClient{stream}
	return x, nil
}

type Services_ChatServiceClient interface {
	Send(*FromClient) error
	Recv() (*FromServer, error)
	grpc.ClientStream
}

type servicesChatServiceClient struct {
	grpc.ClientStream
}

func (x *servicesChatServiceClient) Send(m *FromClient) error {
	return x.ClientStream.SendMsg(m)
}

func (x *servicesChatServiceClient) Recv() (*FromServer, error) {
	m := new(FromServer)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *servicesClient) CommandService(ctx context.Context, in *Command, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/pixalquarks.terminalChatServer.Services/CommandService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) GetClients(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Clients, error) {
	out := new(Clients)
	err := c.cc.Invoke(ctx, "/pixalquarks.terminalChatServer.Services/GetClients", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *servicesClient) VerifyName(ctx context.Context, in *ClientName, opts ...grpc.CallOption) (*ClientNameResponse, error) {
	out := new(ClientNameResponse)
	err := c.cc.Invoke(ctx, "/pixalquarks.terminalChatServer.Services/VerifyName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServicesServer is the server API for Services service.
// All implementations must embed UnimplementedServicesServer
// for forward compatibility
type ServicesServer interface {
	ChatService(Services_ChatServiceServer) error
	CommandService(context.Context, *Command) (*emptypb.Empty, error)
	GetClients(context.Context, *emptypb.Empty) (*Clients, error)
	VerifyName(context.Context, *ClientName) (*ClientNameResponse, error)
	mustEmbedUnimplementedServicesServer()
}

// UnimplementedServicesServer must be embedded to have forward compatible implementations.
type UnimplementedServicesServer struct {
}

func (UnimplementedServicesServer) ChatService(Services_ChatServiceServer) error {
	return status.Errorf(codes.Unimplemented, "method ChatService not implemented")
}
func (UnimplementedServicesServer) CommandService(context.Context, *Command) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommandService not implemented")
}
func (UnimplementedServicesServer) GetClients(context.Context, *emptypb.Empty) (*Clients, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClients not implemented")
}
func (UnimplementedServicesServer) VerifyName(context.Context, *ClientName) (*ClientNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyName not implemented")
}
func (UnimplementedServicesServer) mustEmbedUnimplementedServicesServer() {}

// UnsafeServicesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServicesServer will
// result in compilation errors.
type UnsafeServicesServer interface {
	mustEmbedUnimplementedServicesServer()
}

func RegisterServicesServer(s grpc.ServiceRegistrar, srv ServicesServer) {
	s.RegisterService(&Services_ServiceDesc, srv)
}

func _Services_ChatService_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ServicesServer).ChatService(&servicesChatServiceServer{stream})
}

type Services_ChatServiceServer interface {
	Send(*FromServer) error
	Recv() (*FromClient, error)
	grpc.ServerStream
}

type servicesChatServiceServer struct {
	grpc.ServerStream
}

func (x *servicesChatServiceServer) Send(m *FromServer) error {
	return x.ServerStream.SendMsg(m)
}

func (x *servicesChatServiceServer) Recv() (*FromClient, error) {
	m := new(FromClient)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Services_CommandService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).CommandService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pixalquarks.terminalChatServer.Services/CommandService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).CommandService(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_GetClients_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).GetClients(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pixalquarks.terminalChatServer.Services/GetClients",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).GetClients(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Services_VerifyName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServicesServer).VerifyName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pixalquarks.terminalChatServer.Services/VerifyName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServicesServer).VerifyName(ctx, req.(*ClientName))
	}
	return interceptor(ctx, in, info, handler)
}

// Services_ServiceDesc is the grpc.ServiceDesc for Services service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Services_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pixalquarks.terminalChatServer.Services",
	HandlerType: (*ServicesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CommandService",
			Handler:    _Services_CommandService_Handler,
		},
		{
			MethodName: "GetClients",
			Handler:    _Services_GetClients_Handler,
		},
		{
			MethodName: "VerifyName",
			Handler:    _Services_VerifyName_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ChatService",
			Handler:       _Services_ChatService_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/chat.proto",
}

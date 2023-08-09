// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: pois-service-api.proto

package pb

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

// Podr2ApiClient is the client API for Podr2Api service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type Podr2ApiClient interface {
	RequestGenNodeLoginInfo(ctx context.Context, in *RequestNodeLogin, opts ...grpc.CallOption) (*ResponseNodeLogin, error)
	RequestGenTag(ctx context.Context, in *RequestGenTag, opts ...grpc.CallOption) (*ResponseGenTag, error)
	RequestBatchVerify(ctx context.Context, in *RequestBatchVerify, opts ...grpc.CallOption) (*ResponseBatchVerify, error)
}

type podr2ApiClient struct {
	cc grpc.ClientConnInterface
}

func NewPodr2ApiClient(cc grpc.ClientConnInterface) Podr2ApiClient {
	return &podr2ApiClient{cc}
}

func (c *podr2ApiClient) RequestGenNodeLoginInfo(ctx context.Context, in *RequestNodeLogin, opts ...grpc.CallOption) (*ResponseNodeLogin, error) {
	out := new(ResponseNodeLogin)
	err := c.cc.Invoke(ctx, "/Podr2Api/request_gen_node_login_info", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podr2ApiClient) RequestGenTag(ctx context.Context, in *RequestGenTag, opts ...grpc.CallOption) (*ResponseGenTag, error) {
	out := new(ResponseGenTag)
	err := c.cc.Invoke(ctx, "/Podr2Api/request_gen_tag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *podr2ApiClient) RequestBatchVerify(ctx context.Context, in *RequestBatchVerify, opts ...grpc.CallOption) (*ResponseBatchVerify, error) {
	out := new(ResponseBatchVerify)
	err := c.cc.Invoke(ctx, "/Podr2Api/request_batch_verify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Podr2ApiServer is the server API for Podr2Api service.
// All implementations must embed UnimplementedPodr2ApiServer
// for forward compatibility
type Podr2ApiServer interface {
	RequestGenNodeLoginInfo(context.Context, *RequestNodeLogin) (*ResponseNodeLogin, error)
	RequestGenTag(context.Context, *RequestGenTag) (*ResponseGenTag, error)
	RequestBatchVerify(context.Context, *RequestBatchVerify) (*ResponseBatchVerify, error)
	mustEmbedUnimplementedPodr2ApiServer()
}

// UnimplementedPodr2ApiServer must be embedded to have forward compatible implementations.
type UnimplementedPodr2ApiServer struct {
}

func (UnimplementedPodr2ApiServer) RequestGenNodeLoginInfo(context.Context, *RequestNodeLogin) (*ResponseNodeLogin, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestGenNodeLoginInfo not implemented")
}
func (UnimplementedPodr2ApiServer) RequestGenTag(context.Context, *RequestGenTag) (*ResponseGenTag, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestGenTag not implemented")
}
func (UnimplementedPodr2ApiServer) RequestBatchVerify(context.Context, *RequestBatchVerify) (*ResponseBatchVerify, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestBatchVerify not implemented")
}
func (UnimplementedPodr2ApiServer) mustEmbedUnimplementedPodr2ApiServer() {}

// UnsafePodr2ApiServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to Podr2ApiServer will
// result in compilation errors.
type UnsafePodr2ApiServer interface {
	mustEmbedUnimplementedPodr2ApiServer()
}

func RegisterPodr2ApiServer(s grpc.ServiceRegistrar, srv Podr2ApiServer) {
	s.RegisterService(&Podr2Api_ServiceDesc, srv)
}

func _Podr2Api_RequestGenNodeLoginInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestNodeLogin)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Podr2ApiServer).RequestGenNodeLoginInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Podr2Api/request_gen_node_login_info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Podr2ApiServer).RequestGenNodeLoginInfo(ctx, req.(*RequestNodeLogin))
	}
	return interceptor(ctx, in, info, handler)
}

func _Podr2Api_RequestGenTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestGenTag)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Podr2ApiServer).RequestGenTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Podr2Api/request_gen_tag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Podr2ApiServer).RequestGenTag(ctx, req.(*RequestGenTag))
	}
	return interceptor(ctx, in, info, handler)
}

func _Podr2Api_RequestBatchVerify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestBatchVerify)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Podr2ApiServer).RequestBatchVerify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Podr2Api/request_batch_verify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(Podr2ApiServer).RequestBatchVerify(ctx, req.(*RequestBatchVerify))
	}
	return interceptor(ctx, in, info, handler)
}

// Podr2Api_ServiceDesc is the grpc.ServiceDesc for Podr2Api service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Podr2Api_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Podr2Api",
	HandlerType: (*Podr2ApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "request_gen_node_login_info",
			Handler:    _Podr2Api_RequestGenNodeLoginInfo_Handler,
		},
		{
			MethodName: "request_gen_tag",
			Handler:    _Podr2Api_RequestGenTag_Handler,
		},
		{
			MethodName: "request_batch_verify",
			Handler:    _Podr2Api_RequestBatchVerify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pois-service-api.proto",
}
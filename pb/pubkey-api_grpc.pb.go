// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: pubkey-api.proto

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

const (
	CesealPubkeysProvider_GetIdentityPubkey_FullMethodName = "/ceseal.pubkeys.CesealPubkeysProvider/get_identity_pubkey"
	CesealPubkeysProvider_GetMasterPubkey_FullMethodName   = "/ceseal.pubkeys.CesealPubkeysProvider/get_master_pubkey"
	CesealPubkeysProvider_GetPodr2Pubkey_FullMethodName    = "/ceseal.pubkeys.CesealPubkeysProvider/get_podr2_pubkey"
)

// CesealPubkeysProviderClient is the client API for CesealPubkeysProvider service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CesealPubkeysProviderClient interface {
	// Get the Ceseal identity public key
	GetIdentityPubkey(ctx context.Context, in *Request, opts ...grpc.CallOption) (*IdentityPubkeyResponse, error)
	// Get the master public key
	GetMasterPubkey(ctx context.Context, in *Request, opts ...grpc.CallOption) (*MasterPubkeyResponse, error)
	// Get the PORD2 public key
	GetPodr2Pubkey(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Podr2PubkeyResponse, error)
}

type cesealPubkeysProviderClient struct {
	cc grpc.ClientConnInterface
}

func NewCesealPubkeysProviderClient(cc grpc.ClientConnInterface) CesealPubkeysProviderClient {
	return &cesealPubkeysProviderClient{cc}
}

func (c *cesealPubkeysProviderClient) GetIdentityPubkey(ctx context.Context, in *Request, opts ...grpc.CallOption) (*IdentityPubkeyResponse, error) {
	out := new(IdentityPubkeyResponse)
	err := c.cc.Invoke(ctx, CesealPubkeysProvider_GetIdentityPubkey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cesealPubkeysProviderClient) GetMasterPubkey(ctx context.Context, in *Request, opts ...grpc.CallOption) (*MasterPubkeyResponse, error) {
	out := new(MasterPubkeyResponse)
	err := c.cc.Invoke(ctx, CesealPubkeysProvider_GetMasterPubkey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cesealPubkeysProviderClient) GetPodr2Pubkey(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Podr2PubkeyResponse, error) {
	out := new(Podr2PubkeyResponse)
	err := c.cc.Invoke(ctx, CesealPubkeysProvider_GetPodr2Pubkey_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CesealPubkeysProviderServer is the server API for CesealPubkeysProvider service.
// All implementations must embed UnimplementedCesealPubkeysProviderServer
// for forward compatibility
type CesealPubkeysProviderServer interface {
	// Get the Ceseal identity public key
	GetIdentityPubkey(context.Context, *Request) (*IdentityPubkeyResponse, error)
	// Get the master public key
	GetMasterPubkey(context.Context, *Request) (*MasterPubkeyResponse, error)
	// Get the PORD2 public key
	GetPodr2Pubkey(context.Context, *Request) (*Podr2PubkeyResponse, error)
	mustEmbedUnimplementedCesealPubkeysProviderServer()
}

// UnimplementedCesealPubkeysProviderServer must be embedded to have forward compatible implementations.
type UnimplementedCesealPubkeysProviderServer struct {
}

func (UnimplementedCesealPubkeysProviderServer) GetIdentityPubkey(context.Context, *Request) (*IdentityPubkeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIdentityPubkey not implemented")
}
func (UnimplementedCesealPubkeysProviderServer) GetMasterPubkey(context.Context, *Request) (*MasterPubkeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMasterPubkey not implemented")
}
func (UnimplementedCesealPubkeysProviderServer) GetPodr2Pubkey(context.Context, *Request) (*Podr2PubkeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPodr2Pubkey not implemented")
}
func (UnimplementedCesealPubkeysProviderServer) mustEmbedUnimplementedCesealPubkeysProviderServer() {}

// UnsafeCesealPubkeysProviderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CesealPubkeysProviderServer will
// result in compilation errors.
type UnsafeCesealPubkeysProviderServer interface {
	mustEmbedUnimplementedCesealPubkeysProviderServer()
}

func RegisterCesealPubkeysProviderServer(s grpc.ServiceRegistrar, srv CesealPubkeysProviderServer) {
	s.RegisterService(&CesealPubkeysProvider_ServiceDesc, srv)
}

func _CesealPubkeysProvider_GetIdentityPubkey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CesealPubkeysProviderServer).GetIdentityPubkey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CesealPubkeysProvider_GetIdentityPubkey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CesealPubkeysProviderServer).GetIdentityPubkey(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _CesealPubkeysProvider_GetMasterPubkey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CesealPubkeysProviderServer).GetMasterPubkey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CesealPubkeysProvider_GetMasterPubkey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CesealPubkeysProviderServer).GetMasterPubkey(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _CesealPubkeysProvider_GetPodr2Pubkey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CesealPubkeysProviderServer).GetPodr2Pubkey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CesealPubkeysProvider_GetPodr2Pubkey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CesealPubkeysProviderServer).GetPodr2Pubkey(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// CesealPubkeysProvider_ServiceDesc is the grpc.ServiceDesc for CesealPubkeysProvider service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CesealPubkeysProvider_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ceseal.pubkeys.CesealPubkeysProvider",
	HandlerType: (*CesealPubkeysProviderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "get_identity_pubkey",
			Handler:    _CesealPubkeysProvider_GetIdentityPubkey_Handler,
		},
		{
			MethodName: "get_master_pubkey",
			Handler:    _CesealPubkeysProvider_GetMasterPubkey_Handler,
		},
		{
			MethodName: "get_podr2_pubkey",
			Handler:    _CesealPubkeysProvider_GetPodr2Pubkey_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pubkey-api.proto",
}

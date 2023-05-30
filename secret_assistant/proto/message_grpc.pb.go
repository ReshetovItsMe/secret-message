// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: message.proto

package message

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
	SecretAssistant_Encrypt_FullMethodName = "/main.SecretAssistant/encrypt"
	SecretAssistant_Decrypt_FullMethodName = "/main.SecretAssistant/decrypt"
)

// SecretAssistantClient is the client API for SecretAssistant service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretAssistantClient interface {
	Encrypt(ctx context.Context, in *EncryptRequestMessage, opts ...grpc.CallOption) (*EncryptedMessageResponse, error)
	Decrypt(ctx context.Context, in *DecryptRequestMessage, opts ...grpc.CallOption) (*DecryptedMessageResponse, error)
}

type secretAssistantClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretAssistantClient(cc grpc.ClientConnInterface) SecretAssistantClient {
	return &secretAssistantClient{cc}
}

func (c *secretAssistantClient) Encrypt(ctx context.Context, in *EncryptRequestMessage, opts ...grpc.CallOption) (*EncryptedMessageResponse, error) {
	out := new(EncryptedMessageResponse)
	err := c.cc.Invoke(ctx, SecretAssistant_Encrypt_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretAssistantClient) Decrypt(ctx context.Context, in *DecryptRequestMessage, opts ...grpc.CallOption) (*DecryptedMessageResponse, error) {
	out := new(DecryptedMessageResponse)
	err := c.cc.Invoke(ctx, SecretAssistant_Decrypt_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretAssistantServer is the server API for SecretAssistant service.
// All implementations must embed UnimplementedSecretAssistantServer
// for forward compatibility
type SecretAssistantServer interface {
	Encrypt(context.Context, *EncryptRequestMessage) (*EncryptedMessageResponse, error)
	Decrypt(context.Context, *DecryptRequestMessage) (*DecryptedMessageResponse, error)
	mustEmbedUnimplementedSecretAssistantServer()
}

// UnimplementedSecretAssistantServer must be embedded to have forward compatible implementations.
type UnimplementedSecretAssistantServer struct {
}

func (UnimplementedSecretAssistantServer) Encrypt(context.Context, *EncryptRequestMessage) (*EncryptedMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Encrypt not implemented")
}
func (UnimplementedSecretAssistantServer) Decrypt(context.Context, *DecryptRequestMessage) (*DecryptedMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Decrypt not implemented")
}
func (UnimplementedSecretAssistantServer) mustEmbedUnimplementedSecretAssistantServer() {}

// UnsafeSecretAssistantServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretAssistantServer will
// result in compilation errors.
type UnsafeSecretAssistantServer interface {
	mustEmbedUnimplementedSecretAssistantServer()
}

func RegisterSecretAssistantServer(s grpc.ServiceRegistrar, srv SecretAssistantServer) {
	s.RegisterService(&SecretAssistant_ServiceDesc, srv)
}

func _SecretAssistant_Encrypt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EncryptRequestMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretAssistantServer).Encrypt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SecretAssistant_Encrypt_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretAssistantServer).Encrypt(ctx, req.(*EncryptRequestMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretAssistant_Decrypt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DecryptRequestMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretAssistantServer).Decrypt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SecretAssistant_Decrypt_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretAssistantServer).Decrypt(ctx, req.(*DecryptRequestMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// SecretAssistant_ServiceDesc is the grpc.ServiceDesc for SecretAssistant service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SecretAssistant_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.SecretAssistant",
	HandlerType: (*SecretAssistantServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "encrypt",
			Handler:    _SecretAssistant_Encrypt_Handler,
		},
		{
			MethodName: "decrypt",
			Handler:    _SecretAssistant_Decrypt_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "message.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/v1/word.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	WordService_Get_FullMethodName  = "/proto.v1.WordService/Get"
	WordService_List_FullMethodName = "/proto.v1.WordService/List"
)

// WordServiceClient is the client API for WordService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// WordService 提供单词相关的服务
type WordServiceClient interface {
	// Get 获取单词详情
	Get(ctx context.Context, in *WordServiceGetRequest, opts ...grpc.CallOption) (*WordServiceGetResponse, error)
	// List 获取单词列表
	List(ctx context.Context, in *WordServiceListRequest, opts ...grpc.CallOption) (*WordServiceListResponse, error)
}

type wordServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWordServiceClient(cc grpc.ClientConnInterface) WordServiceClient {
	return &wordServiceClient{cc}
}

func (c *wordServiceClient) Get(ctx context.Context, in *WordServiceGetRequest, opts ...grpc.CallOption) (*WordServiceGetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WordServiceGetResponse)
	err := c.cc.Invoke(ctx, WordService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *wordServiceClient) List(ctx context.Context, in *WordServiceListRequest, opts ...grpc.CallOption) (*WordServiceListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WordServiceListResponse)
	err := c.cc.Invoke(ctx, WordService_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WordServiceServer is the server API for WordService service.
// All implementations must embed UnimplementedWordServiceServer
// for forward compatibility.
//
// WordService 提供单词相关的服务
type WordServiceServer interface {
	// Get 获取单词详情
	Get(context.Context, *WordServiceGetRequest) (*WordServiceGetResponse, error)
	// List 获取单词列表
	List(context.Context, *WordServiceListRequest) (*WordServiceListResponse, error)
	mustEmbedUnimplementedWordServiceServer()
}

// UnimplementedWordServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedWordServiceServer struct{}

func (UnimplementedWordServiceServer) Get(context.Context, *WordServiceGetRequest) (*WordServiceGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedWordServiceServer) List(context.Context, *WordServiceListRequest) (*WordServiceListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedWordServiceServer) mustEmbedUnimplementedWordServiceServer() {}
func (UnimplementedWordServiceServer) testEmbeddedByValue()                     {}

// UnsafeWordServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WordServiceServer will
// result in compilation errors.
type UnsafeWordServiceServer interface {
	mustEmbedUnimplementedWordServiceServer()
}

func RegisterWordServiceServer(s grpc.ServiceRegistrar, srv WordServiceServer) {
	// If the following call pancis, it indicates UnimplementedWordServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&WordService_ServiceDesc, srv)
}

func _WordService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WordServiceGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WordServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WordService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WordServiceServer).Get(ctx, req.(*WordServiceGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WordService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WordServiceListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WordServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WordService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WordServiceServer).List(ctx, req.(*WordServiceListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WordService_ServiceDesc is the grpc.ServiceDesc for WordService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WordService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.v1.WordService",
	HandlerType: (*WordServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _WordService_Get_Handler,
		},
		{
			MethodName: "List",
			Handler:    _WordService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/word.proto",
}

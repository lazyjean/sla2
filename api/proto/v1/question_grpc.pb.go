// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: proto/v1/question.proto

// ============= 基础定义 =============

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
	QuestionService_Get_FullMethodName     = "/proto.v1.QuestionService/Get"
	QuestionService_Create_FullMethodName  = "/proto.v1.QuestionService/Create"
	QuestionService_Search_FullMethodName  = "/proto.v1.QuestionService/Search"
	QuestionService_Update_FullMethodName  = "/proto.v1.QuestionService/Update"
	QuestionService_Publish_FullMethodName = "/proto.v1.QuestionService/Publish"
	QuestionService_Delete_FullMethodName  = "/proto.v1.QuestionService/Delete"
)

// QuestionServiceClient is the client API for QuestionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QuestionServiceClient interface {
	Get(ctx context.Context, in *QuestionServiceGetRequest, opts ...grpc.CallOption) (*QuestionServiceGetResponse, error)
	Create(ctx context.Context, in *QuestionServiceCreateRequest, opts ...grpc.CallOption) (*QuestionServiceCreateResponse, error)
	Search(ctx context.Context, in *QuestionServiceSearchRequest, opts ...grpc.CallOption) (*QuestionServiceSearchResponse, error)
	Update(ctx context.Context, in *QuestionServiceUpdateRequest, opts ...grpc.CallOption) (*QuestionServiceUpdateResponse, error)
	Publish(ctx context.Context, in *QuestionServicePublishRequest, opts ...grpc.CallOption) (*QuestionServicePublishResponse, error)
	Delete(ctx context.Context, in *QuestionServiceDeleteRequest, opts ...grpc.CallOption) (*QuestionServiceDeleteResponse, error)
}

type questionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQuestionServiceClient(cc grpc.ClientConnInterface) QuestionServiceClient {
	return &questionServiceClient{cc}
}

func (c *questionServiceClient) Get(ctx context.Context, in *QuestionServiceGetRequest, opts ...grpc.CallOption) (*QuestionServiceGetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuestionServiceGetResponse)
	err := c.cc.Invoke(ctx, QuestionService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionServiceClient) Create(ctx context.Context, in *QuestionServiceCreateRequest, opts ...grpc.CallOption) (*QuestionServiceCreateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuestionServiceCreateResponse)
	err := c.cc.Invoke(ctx, QuestionService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionServiceClient) Search(ctx context.Context, in *QuestionServiceSearchRequest, opts ...grpc.CallOption) (*QuestionServiceSearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuestionServiceSearchResponse)
	err := c.cc.Invoke(ctx, QuestionService_Search_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionServiceClient) Update(ctx context.Context, in *QuestionServiceUpdateRequest, opts ...grpc.CallOption) (*QuestionServiceUpdateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuestionServiceUpdateResponse)
	err := c.cc.Invoke(ctx, QuestionService_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionServiceClient) Publish(ctx context.Context, in *QuestionServicePublishRequest, opts ...grpc.CallOption) (*QuestionServicePublishResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuestionServicePublishResponse)
	err := c.cc.Invoke(ctx, QuestionService_Publish_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionServiceClient) Delete(ctx context.Context, in *QuestionServiceDeleteRequest, opts ...grpc.CallOption) (*QuestionServiceDeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuestionServiceDeleteResponse)
	err := c.cc.Invoke(ctx, QuestionService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QuestionServiceServer is the server API for QuestionService service.
// All implementations must embed UnimplementedQuestionServiceServer
// for forward compatibility.
type QuestionServiceServer interface {
	Get(context.Context, *QuestionServiceGetRequest) (*QuestionServiceGetResponse, error)
	Create(context.Context, *QuestionServiceCreateRequest) (*QuestionServiceCreateResponse, error)
	Search(context.Context, *QuestionServiceSearchRequest) (*QuestionServiceSearchResponse, error)
	Update(context.Context, *QuestionServiceUpdateRequest) (*QuestionServiceUpdateResponse, error)
	Publish(context.Context, *QuestionServicePublishRequest) (*QuestionServicePublishResponse, error)
	Delete(context.Context, *QuestionServiceDeleteRequest) (*QuestionServiceDeleteResponse, error)
	mustEmbedUnimplementedQuestionServiceServer()
}

// UnimplementedQuestionServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQuestionServiceServer struct{}

func (UnimplementedQuestionServiceServer) Get(context.Context, *QuestionServiceGetRequest) (*QuestionServiceGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedQuestionServiceServer) Create(context.Context, *QuestionServiceCreateRequest) (*QuestionServiceCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedQuestionServiceServer) Search(context.Context, *QuestionServiceSearchRequest) (*QuestionServiceSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedQuestionServiceServer) Update(context.Context, *QuestionServiceUpdateRequest) (*QuestionServiceUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedQuestionServiceServer) Publish(context.Context, *QuestionServicePublishRequest) (*QuestionServicePublishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedQuestionServiceServer) Delete(context.Context, *QuestionServiceDeleteRequest) (*QuestionServiceDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedQuestionServiceServer) mustEmbedUnimplementedQuestionServiceServer() {}
func (UnimplementedQuestionServiceServer) testEmbeddedByValue()                         {}

// UnsafeQuestionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QuestionServiceServer will
// result in compilation errors.
type UnsafeQuestionServiceServer interface {
	mustEmbedUnimplementedQuestionServiceServer()
}

func RegisterQuestionServiceServer(s grpc.ServiceRegistrar, srv QuestionServiceServer) {
	// If the following call pancis, it indicates UnimplementedQuestionServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&QuestionService_ServiceDesc, srv)
}

func _QuestionService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestionServiceGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QuestionService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionServiceServer).Get(ctx, req.(*QuestionServiceGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestionServiceCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QuestionService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionServiceServer).Create(ctx, req.(*QuestionServiceCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestionServiceSearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QuestionService_Search_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionServiceServer).Search(ctx, req.(*QuestionServiceSearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestionServiceUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QuestionService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionServiceServer).Update(ctx, req.(*QuestionServiceUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionService_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestionServicePublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionServiceServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QuestionService_Publish_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionServiceServer).Publish(ctx, req.(*QuestionServicePublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestionServiceDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QuestionService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionServiceServer).Delete(ctx, req.(*QuestionServiceDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// QuestionService_ServiceDesc is the grpc.ServiceDesc for QuestionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QuestionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.v1.QuestionService",
	HandlerType: (*QuestionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _QuestionService_Get_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _QuestionService_Create_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _QuestionService_Search_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _QuestionService_Update_Handler,
		},
		{
			MethodName: "Publish",
			Handler:    _QuestionService_Publish_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _QuestionService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/question.proto",
}

const (
	TagService_ListTag_FullMethodName   = "/proto.v1.TagService/ListTag"
	TagService_CreateTag_FullMethodName = "/proto.v1.TagService/CreateTag"
	TagService_UpdateTag_FullMethodName = "/proto.v1.TagService/UpdateTag"
	TagService_DeleteTag_FullMethodName = "/proto.v1.TagService/DeleteTag"
)

// TagServiceClient is the client API for TagService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TagServiceClient interface {
	ListTag(ctx context.Context, in *ListTagRequest, opts ...grpc.CallOption) (*ListTagResponse, error)
	CreateTag(ctx context.Context, in *CreateTagRequest, opts ...grpc.CallOption) (*CreateTagResponse, error)
	UpdateTag(ctx context.Context, in *UpdateTagRequest, opts ...grpc.CallOption) (*UpdateTagResponse, error)
	// delete tag for test
	DeleteTag(ctx context.Context, in *DeleteTagRequest, opts ...grpc.CallOption) (*DeleteTagResponse, error)
}

type tagServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTagServiceClient(cc grpc.ClientConnInterface) TagServiceClient {
	return &tagServiceClient{cc}
}

func (c *tagServiceClient) ListTag(ctx context.Context, in *ListTagRequest, opts ...grpc.CallOption) (*ListTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListTagResponse)
	err := c.cc.Invoke(ctx, TagService_ListTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tagServiceClient) CreateTag(ctx context.Context, in *CreateTagRequest, opts ...grpc.CallOption) (*CreateTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateTagResponse)
	err := c.cc.Invoke(ctx, TagService_CreateTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tagServiceClient) UpdateTag(ctx context.Context, in *UpdateTagRequest, opts ...grpc.CallOption) (*UpdateTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateTagResponse)
	err := c.cc.Invoke(ctx, TagService_UpdateTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tagServiceClient) DeleteTag(ctx context.Context, in *DeleteTagRequest, opts ...grpc.CallOption) (*DeleteTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteTagResponse)
	err := c.cc.Invoke(ctx, TagService_DeleteTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TagServiceServer is the server API for TagService service.
// All implementations must embed UnimplementedTagServiceServer
// for forward compatibility.
type TagServiceServer interface {
	ListTag(context.Context, *ListTagRequest) (*ListTagResponse, error)
	CreateTag(context.Context, *CreateTagRequest) (*CreateTagResponse, error)
	UpdateTag(context.Context, *UpdateTagRequest) (*UpdateTagResponse, error)
	// delete tag for test
	DeleteTag(context.Context, *DeleteTagRequest) (*DeleteTagResponse, error)
	mustEmbedUnimplementedTagServiceServer()
}

// UnimplementedTagServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTagServiceServer struct{}

func (UnimplementedTagServiceServer) ListTag(context.Context, *ListTagRequest) (*ListTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTag not implemented")
}
func (UnimplementedTagServiceServer) CreateTag(context.Context, *CreateTagRequest) (*CreateTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTag not implemented")
}
func (UnimplementedTagServiceServer) UpdateTag(context.Context, *UpdateTagRequest) (*UpdateTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTag not implemented")
}
func (UnimplementedTagServiceServer) DeleteTag(context.Context, *DeleteTagRequest) (*DeleteTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTag not implemented")
}
func (UnimplementedTagServiceServer) mustEmbedUnimplementedTagServiceServer() {}
func (UnimplementedTagServiceServer) testEmbeddedByValue()                    {}

// UnsafeTagServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TagServiceServer will
// result in compilation errors.
type UnsafeTagServiceServer interface {
	mustEmbedUnimplementedTagServiceServer()
}

func RegisterTagServiceServer(s grpc.ServiceRegistrar, srv TagServiceServer) {
	// If the following call pancis, it indicates UnimplementedTagServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TagService_ServiceDesc, srv)
}

func _TagService_ListTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TagServiceServer).ListTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TagService_ListTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TagServiceServer).ListTag(ctx, req.(*ListTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TagService_CreateTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TagServiceServer).CreateTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TagService_CreateTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TagServiceServer).CreateTag(ctx, req.(*CreateTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TagService_UpdateTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TagServiceServer).UpdateTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TagService_UpdateTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TagServiceServer).UpdateTag(ctx, req.(*UpdateTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TagService_DeleteTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TagServiceServer).DeleteTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TagService_DeleteTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TagServiceServer).DeleteTag(ctx, req.(*DeleteTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TagService_ServiceDesc is the grpc.ServiceDesc for TagService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TagService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.v1.TagService",
	HandlerType: (*TagServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListTag",
			Handler:    _TagService_ListTag_Handler,
		},
		{
			MethodName: "CreateTag",
			Handler:    _TagService_CreateTag_Handler,
		},
		{
			MethodName: "UpdateTag",
			Handler:    _TagService_UpdateTag_Handler,
		},
		{
			MethodName: "DeleteTag",
			Handler:    _TagService_DeleteTag_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/question.proto",
}

const (
	QuestionTagService_ListTag_FullMethodName   = "/proto.v1.QuestionTagService/ListTag"
	QuestionTagService_CreateTag_FullMethodName = "/proto.v1.QuestionTagService/CreateTag"
	QuestionTagService_UpdateTag_FullMethodName = "/proto.v1.QuestionTagService/UpdateTag"
	QuestionTagService_DeleteTag_FullMethodName = "/proto.v1.QuestionTagService/DeleteTag"
)

// QuestionTagServiceClient is the client API for QuestionTagService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QuestionTagServiceClient interface {
	ListTag(ctx context.Context, in *QuestionTagServiceListTagRequest, opts ...grpc.CallOption) (*QuestionTagServiceListTagResponse, error)
	CreateTag(ctx context.Context, in *QuestionTagServiceCreateTagRequest, opts ...grpc.CallOption) (*QuestionTagServiceCreateTagResponse, error)
	UpdateTag(ctx context.Context, in *QuestionTagServiceUpdateTagRequest, opts ...grpc.CallOption) (*QuestionTagServiceUpdateTagResponse, error)
	DeleteTag(ctx context.Context, in *QuestionTagServiceDeleteTagRequest, opts ...grpc.CallOption) (*QuestionTagServiceDeleteTagResponse, error)
}

type questionTagServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQuestionTagServiceClient(cc grpc.ClientConnInterface) QuestionTagServiceClient {
	return &questionTagServiceClient{cc}
}

func (c *questionTagServiceClient) ListTag(ctx context.Context, in *QuestionTagServiceListTagRequest, opts ...grpc.CallOption) (*QuestionTagServiceListTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuestionTagServiceListTagResponse)
	err := c.cc.Invoke(ctx, QuestionTagService_ListTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionTagServiceClient) CreateTag(ctx context.Context, in *QuestionTagServiceCreateTagRequest, opts ...grpc.CallOption) (*QuestionTagServiceCreateTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuestionTagServiceCreateTagResponse)
	err := c.cc.Invoke(ctx, QuestionTagService_CreateTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionTagServiceClient) UpdateTag(ctx context.Context, in *QuestionTagServiceUpdateTagRequest, opts ...grpc.CallOption) (*QuestionTagServiceUpdateTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuestionTagServiceUpdateTagResponse)
	err := c.cc.Invoke(ctx, QuestionTagService_UpdateTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *questionTagServiceClient) DeleteTag(ctx context.Context, in *QuestionTagServiceDeleteTagRequest, opts ...grpc.CallOption) (*QuestionTagServiceDeleteTagResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QuestionTagServiceDeleteTagResponse)
	err := c.cc.Invoke(ctx, QuestionTagService_DeleteTag_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QuestionTagServiceServer is the server API for QuestionTagService service.
// All implementations must embed UnimplementedQuestionTagServiceServer
// for forward compatibility.
type QuestionTagServiceServer interface {
	ListTag(context.Context, *QuestionTagServiceListTagRequest) (*QuestionTagServiceListTagResponse, error)
	CreateTag(context.Context, *QuestionTagServiceCreateTagRequest) (*QuestionTagServiceCreateTagResponse, error)
	UpdateTag(context.Context, *QuestionTagServiceUpdateTagRequest) (*QuestionTagServiceUpdateTagResponse, error)
	DeleteTag(context.Context, *QuestionTagServiceDeleteTagRequest) (*QuestionTagServiceDeleteTagResponse, error)
	mustEmbedUnimplementedQuestionTagServiceServer()
}

// UnimplementedQuestionTagServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQuestionTagServiceServer struct{}

func (UnimplementedQuestionTagServiceServer) ListTag(context.Context, *QuestionTagServiceListTagRequest) (*QuestionTagServiceListTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTag not implemented")
}
func (UnimplementedQuestionTagServiceServer) CreateTag(context.Context, *QuestionTagServiceCreateTagRequest) (*QuestionTagServiceCreateTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTag not implemented")
}
func (UnimplementedQuestionTagServiceServer) UpdateTag(context.Context, *QuestionTagServiceUpdateTagRequest) (*QuestionTagServiceUpdateTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTag not implemented")
}
func (UnimplementedQuestionTagServiceServer) DeleteTag(context.Context, *QuestionTagServiceDeleteTagRequest) (*QuestionTagServiceDeleteTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTag not implemented")
}
func (UnimplementedQuestionTagServiceServer) mustEmbedUnimplementedQuestionTagServiceServer() {}
func (UnimplementedQuestionTagServiceServer) testEmbeddedByValue()                            {}

// UnsafeQuestionTagServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QuestionTagServiceServer will
// result in compilation errors.
type UnsafeQuestionTagServiceServer interface {
	mustEmbedUnimplementedQuestionTagServiceServer()
}

func RegisterQuestionTagServiceServer(s grpc.ServiceRegistrar, srv QuestionTagServiceServer) {
	// If the following call pancis, it indicates UnimplementedQuestionTagServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&QuestionTagService_ServiceDesc, srv)
}

func _QuestionTagService_ListTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestionTagServiceListTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionTagServiceServer).ListTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QuestionTagService_ListTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionTagServiceServer).ListTag(ctx, req.(*QuestionTagServiceListTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionTagService_CreateTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestionTagServiceCreateTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionTagServiceServer).CreateTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QuestionTagService_CreateTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionTagServiceServer).CreateTag(ctx, req.(*QuestionTagServiceCreateTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionTagService_UpdateTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestionTagServiceUpdateTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionTagServiceServer).UpdateTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QuestionTagService_UpdateTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionTagServiceServer).UpdateTag(ctx, req.(*QuestionTagServiceUpdateTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuestionTagService_DeleteTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuestionTagServiceDeleteTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuestionTagServiceServer).DeleteTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QuestionTagService_DeleteTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuestionTagServiceServer).DeleteTag(ctx, req.(*QuestionTagServiceDeleteTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// QuestionTagService_ServiceDesc is the grpc.ServiceDesc for QuestionTagService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QuestionTagService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.v1.QuestionTagService",
	HandlerType: (*QuestionTagServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListTag",
			Handler:    _QuestionTagService_ListTag_Handler,
		},
		{
			MethodName: "CreateTag",
			Handler:    _QuestionTagService_CreateTag_Handler,
		},
		{
			MethodName: "UpdateTag",
			Handler:    _QuestionTagService_UpdateTag_Handler,
		},
		{
			MethodName: "DeleteTag",
			Handler:    _QuestionTagService_DeleteTag_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/question.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/v1/course.proto

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
	CourseService_Create_FullMethodName        = "/proto.v1.CourseService/Create"
	CourseService_Update_FullMethodName        = "/proto.v1.CourseService/Update"
	CourseService_Get_FullMethodName           = "/proto.v1.CourseService/Get"
	CourseService_List_FullMethodName          = "/proto.v1.CourseService/List"
	CourseService_Delete_FullMethodName        = "/proto.v1.CourseService/Delete"
	CourseService_Search_FullMethodName        = "/proto.v1.CourseService/Search"
	CourseService_CreateSection_FullMethodName = "/proto.v1.CourseService/CreateSection"
	CourseService_UpdateSection_FullMethodName = "/proto.v1.CourseService/UpdateSection"
	CourseService_DeleteSection_FullMethodName = "/proto.v1.CourseService/DeleteSection"
	CourseService_CreateUnit_FullMethodName    = "/proto.v1.CourseService/CreateUnit"
	CourseService_UpdateUnit_FullMethodName    = "/proto.v1.CourseService/UpdateUnit"
	CourseService_DeleteUnit_FullMethodName    = "/proto.v1.CourseService/DeleteUnit"
)

// CourseServiceClient is the client API for CourseService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 由于课程是相对比较固定的内容, 所以我们这里设计为一次获取课程所有结构数据, 没必要分多次获取
// CourseService 提供课程相关的服务
type CourseServiceClient interface {
	// Create 创建课程
	Create(ctx context.Context, in *CourseServiceCreateRequest, opts ...grpc.CallOption) (*CourseServiceCreateResponse, error)
	// Update 更新课程
	Update(ctx context.Context, in *CourseServiceUpdateRequest, opts ...grpc.CallOption) (*CourseServiceUpdateResponse, error)
	// Get 获取课程详情
	Get(ctx context.Context, in *CourseServiceGetRequest, opts ...grpc.CallOption) (*CourseServiceGetResponse, error)
	// List 获取课程列表
	List(ctx context.Context, in *CourseServiceListRequest, opts ...grpc.CallOption) (*CourseServiceListResponse, error)
	// Delete 删除课程
	Delete(ctx context.Context, in *CourseServiceDeleteRequest, opts ...grpc.CallOption) (*CourseServiceDeleteResponse, error)
	// Search 搜索课程
	Search(ctx context.Context, in *CourseServiceSearchRequest, opts ...grpc.CallOption) (*CourseServiceSearchResponse, error)
	// CreateSection 创建课程章节
	CreateSection(ctx context.Context, in *CourseServiceCreateSectionRequest, opts ...grpc.CallOption) (*CourseServiceCreateSectionResponse, error)
	// UpdateSection 更新课程章节
	UpdateSection(ctx context.Context, in *CourseServiceUpdateSectionRequest, opts ...grpc.CallOption) (*CourseServiceUpdateSectionResponse, error)
	// DeleteSection 删除课程章节
	DeleteSection(ctx context.Context, in *CourseServiceDeleteSectionRequest, opts ...grpc.CallOption) (*CourseServiceDeleteSectionResponse, error)
	// CreateUnit 创建章节单元
	CreateUnit(ctx context.Context, in *CourseServiceCreateUnitRequest, opts ...grpc.CallOption) (*CourseServiceCreateUnitResponse, error)
	// UpdateUnit 更新章节单元
	UpdateUnit(ctx context.Context, in *CourseServiceUpdateUnitRequest, opts ...grpc.CallOption) (*CourseServiceUpdateUnitResponse, error)
	// DeleteUnit 删除章节单元
	DeleteUnit(ctx context.Context, in *CourseServiceDeleteUnitRequest, opts ...grpc.CallOption) (*CourseServiceDeleteUnitResponse, error)
}

type courseServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCourseServiceClient(cc grpc.ClientConnInterface) CourseServiceClient {
	return &courseServiceClient{cc}
}

func (c *courseServiceClient) Create(ctx context.Context, in *CourseServiceCreateRequest, opts ...grpc.CallOption) (*CourseServiceCreateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceCreateResponse)
	err := c.cc.Invoke(ctx, CourseService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) Update(ctx context.Context, in *CourseServiceUpdateRequest, opts ...grpc.CallOption) (*CourseServiceUpdateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceUpdateResponse)
	err := c.cc.Invoke(ctx, CourseService_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) Get(ctx context.Context, in *CourseServiceGetRequest, opts ...grpc.CallOption) (*CourseServiceGetResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceGetResponse)
	err := c.cc.Invoke(ctx, CourseService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) List(ctx context.Context, in *CourseServiceListRequest, opts ...grpc.CallOption) (*CourseServiceListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceListResponse)
	err := c.cc.Invoke(ctx, CourseService_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) Delete(ctx context.Context, in *CourseServiceDeleteRequest, opts ...grpc.CallOption) (*CourseServiceDeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceDeleteResponse)
	err := c.cc.Invoke(ctx, CourseService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) Search(ctx context.Context, in *CourseServiceSearchRequest, opts ...grpc.CallOption) (*CourseServiceSearchResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceSearchResponse)
	err := c.cc.Invoke(ctx, CourseService_Search_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) CreateSection(ctx context.Context, in *CourseServiceCreateSectionRequest, opts ...grpc.CallOption) (*CourseServiceCreateSectionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceCreateSectionResponse)
	err := c.cc.Invoke(ctx, CourseService_CreateSection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) UpdateSection(ctx context.Context, in *CourseServiceUpdateSectionRequest, opts ...grpc.CallOption) (*CourseServiceUpdateSectionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceUpdateSectionResponse)
	err := c.cc.Invoke(ctx, CourseService_UpdateSection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) DeleteSection(ctx context.Context, in *CourseServiceDeleteSectionRequest, opts ...grpc.CallOption) (*CourseServiceDeleteSectionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceDeleteSectionResponse)
	err := c.cc.Invoke(ctx, CourseService_DeleteSection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) CreateUnit(ctx context.Context, in *CourseServiceCreateUnitRequest, opts ...grpc.CallOption) (*CourseServiceCreateUnitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceCreateUnitResponse)
	err := c.cc.Invoke(ctx, CourseService_CreateUnit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) UpdateUnit(ctx context.Context, in *CourseServiceUpdateUnitRequest, opts ...grpc.CallOption) (*CourseServiceUpdateUnitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceUpdateUnitResponse)
	err := c.cc.Invoke(ctx, CourseService_UpdateUnit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *courseServiceClient) DeleteUnit(ctx context.Context, in *CourseServiceDeleteUnitRequest, opts ...grpc.CallOption) (*CourseServiceDeleteUnitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CourseServiceDeleteUnitResponse)
	err := c.cc.Invoke(ctx, CourseService_DeleteUnit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CourseServiceServer is the server API for CourseService service.
// All implementations must embed UnimplementedCourseServiceServer
// for forward compatibility.
//
// 由于课程是相对比较固定的内容, 所以我们这里设计为一次获取课程所有结构数据, 没必要分多次获取
// CourseService 提供课程相关的服务
type CourseServiceServer interface {
	// Create 创建课程
	Create(context.Context, *CourseServiceCreateRequest) (*CourseServiceCreateResponse, error)
	// Update 更新课程
	Update(context.Context, *CourseServiceUpdateRequest) (*CourseServiceUpdateResponse, error)
	// Get 获取课程详情
	Get(context.Context, *CourseServiceGetRequest) (*CourseServiceGetResponse, error)
	// List 获取课程列表
	List(context.Context, *CourseServiceListRequest) (*CourseServiceListResponse, error)
	// Delete 删除课程
	Delete(context.Context, *CourseServiceDeleteRequest) (*CourseServiceDeleteResponse, error)
	// Search 搜索课程
	Search(context.Context, *CourseServiceSearchRequest) (*CourseServiceSearchResponse, error)
	// CreateSection 创建课程章节
	CreateSection(context.Context, *CourseServiceCreateSectionRequest) (*CourseServiceCreateSectionResponse, error)
	// UpdateSection 更新课程章节
	UpdateSection(context.Context, *CourseServiceUpdateSectionRequest) (*CourseServiceUpdateSectionResponse, error)
	// DeleteSection 删除课程章节
	DeleteSection(context.Context, *CourseServiceDeleteSectionRequest) (*CourseServiceDeleteSectionResponse, error)
	// CreateUnit 创建章节单元
	CreateUnit(context.Context, *CourseServiceCreateUnitRequest) (*CourseServiceCreateUnitResponse, error)
	// UpdateUnit 更新章节单元
	UpdateUnit(context.Context, *CourseServiceUpdateUnitRequest) (*CourseServiceUpdateUnitResponse, error)
	// DeleteUnit 删除章节单元
	DeleteUnit(context.Context, *CourseServiceDeleteUnitRequest) (*CourseServiceDeleteUnitResponse, error)
	mustEmbedUnimplementedCourseServiceServer()
}

// UnimplementedCourseServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCourseServiceServer struct{}

func (UnimplementedCourseServiceServer) Create(context.Context, *CourseServiceCreateRequest) (*CourseServiceCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedCourseServiceServer) Update(context.Context, *CourseServiceUpdateRequest) (*CourseServiceUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedCourseServiceServer) Get(context.Context, *CourseServiceGetRequest) (*CourseServiceGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedCourseServiceServer) List(context.Context, *CourseServiceListRequest) (*CourseServiceListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedCourseServiceServer) Delete(context.Context, *CourseServiceDeleteRequest) (*CourseServiceDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedCourseServiceServer) Search(context.Context, *CourseServiceSearchRequest) (*CourseServiceSearchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (UnimplementedCourseServiceServer) CreateSection(context.Context, *CourseServiceCreateSectionRequest) (*CourseServiceCreateSectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSection not implemented")
}
func (UnimplementedCourseServiceServer) UpdateSection(context.Context, *CourseServiceUpdateSectionRequest) (*CourseServiceUpdateSectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateSection not implemented")
}
func (UnimplementedCourseServiceServer) DeleteSection(context.Context, *CourseServiceDeleteSectionRequest) (*CourseServiceDeleteSectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSection not implemented")
}
func (UnimplementedCourseServiceServer) CreateUnit(context.Context, *CourseServiceCreateUnitRequest) (*CourseServiceCreateUnitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUnit not implemented")
}
func (UnimplementedCourseServiceServer) UpdateUnit(context.Context, *CourseServiceUpdateUnitRequest) (*CourseServiceUpdateUnitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUnit not implemented")
}
func (UnimplementedCourseServiceServer) DeleteUnit(context.Context, *CourseServiceDeleteUnitRequest) (*CourseServiceDeleteUnitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUnit not implemented")
}
func (UnimplementedCourseServiceServer) mustEmbedUnimplementedCourseServiceServer() {}
func (UnimplementedCourseServiceServer) testEmbeddedByValue()                       {}

// UnsafeCourseServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CourseServiceServer will
// result in compilation errors.
type UnsafeCourseServiceServer interface {
	mustEmbedUnimplementedCourseServiceServer()
}

func RegisterCourseServiceServer(s grpc.ServiceRegistrar, srv CourseServiceServer) {
	// If the following call pancis, it indicates UnimplementedCourseServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CourseService_ServiceDesc, srv)
}

func _CourseService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).Create(ctx, req.(*CourseServiceCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).Update(ctx, req.(*CourseServiceUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).Get(ctx, req.(*CourseServiceGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).List(ctx, req.(*CourseServiceListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).Delete(ctx, req.(*CourseServiceDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceSearchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_Search_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).Search(ctx, req.(*CourseServiceSearchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_CreateSection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceCreateSectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).CreateSection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_CreateSection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).CreateSection(ctx, req.(*CourseServiceCreateSectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_UpdateSection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceUpdateSectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).UpdateSection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_UpdateSection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).UpdateSection(ctx, req.(*CourseServiceUpdateSectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_DeleteSection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceDeleteSectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).DeleteSection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_DeleteSection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).DeleteSection(ctx, req.(*CourseServiceDeleteSectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_CreateUnit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceCreateUnitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).CreateUnit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_CreateUnit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).CreateUnit(ctx, req.(*CourseServiceCreateUnitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_UpdateUnit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceUpdateUnitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).UpdateUnit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_UpdateUnit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).UpdateUnit(ctx, req.(*CourseServiceUpdateUnitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CourseService_DeleteUnit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CourseServiceDeleteUnitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CourseServiceServer).DeleteUnit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CourseService_DeleteUnit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CourseServiceServer).DeleteUnit(ctx, req.(*CourseServiceDeleteUnitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CourseService_ServiceDesc is the grpc.ServiceDesc for CourseService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CourseService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.v1.CourseService",
	HandlerType: (*CourseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _CourseService_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _CourseService_Update_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _CourseService_Get_Handler,
		},
		{
			MethodName: "List",
			Handler:    _CourseService_List_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _CourseService_Delete_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _CourseService_Search_Handler,
		},
		{
			MethodName: "CreateSection",
			Handler:    _CourseService_CreateSection_Handler,
		},
		{
			MethodName: "UpdateSection",
			Handler:    _CourseService_UpdateSection_Handler,
		},
		{
			MethodName: "DeleteSection",
			Handler:    _CourseService_DeleteSection_Handler,
		},
		{
			MethodName: "CreateUnit",
			Handler:    _CourseService_CreateUnit_Handler,
		},
		{
			MethodName: "UpdateUnit",
			Handler:    _CourseService_UpdateUnit_Handler,
		},
		{
			MethodName: "DeleteUnit",
			Handler:    _CourseService_DeleteUnit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/course.proto",
}

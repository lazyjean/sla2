// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: proto/v1/admin.proto

package pb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 系统状态响应
type AdminServiceCheckSystemStatusResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Initialized   bool                   `protobuf:"varint,1,opt,name=initialized,proto3" json:"initialized,omitempty"` // 系统是否已初始化
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminServiceCheckSystemStatusResponse) Reset() {
	*x = AdminServiceCheckSystemStatusResponse{}
	mi := &file_proto_v1_admin_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminServiceCheckSystemStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminServiceCheckSystemStatusResponse) ProtoMessage() {}

func (x *AdminServiceCheckSystemStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminServiceCheckSystemStatusResponse.ProtoReflect.Descriptor instead.
func (*AdminServiceCheckSystemStatusResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{0}
}

func (x *AdminServiceCheckSystemStatusResponse) GetInitialized() bool {
	if x != nil {
		return x.Initialized
	}
	return false
}

// 获取系统状态请求
type AdminServiceCheckSystemStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminServiceCheckSystemStatusRequest) Reset() {
	*x = AdminServiceCheckSystemStatusRequest{}
	mi := &file_proto_v1_admin_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminServiceCheckSystemStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminServiceCheckSystemStatusRequest) ProtoMessage() {}

func (x *AdminServiceCheckSystemStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminServiceCheckSystemStatusRequest.ProtoReflect.Descriptor instead.
func (*AdminServiceCheckSystemStatusRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{1}
}

// 获取当前管理员请求
type AdminServiceGetCurrentAdminInfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminServiceGetCurrentAdminInfoRequest) Reset() {
	*x = AdminServiceGetCurrentAdminInfoRequest{}
	mi := &file_proto_v1_admin_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminServiceGetCurrentAdminInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminServiceGetCurrentAdminInfoRequest) ProtoMessage() {}

func (x *AdminServiceGetCurrentAdminInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminServiceGetCurrentAdminInfoRequest.ProtoReflect.Descriptor instead.
func (*AdminServiceGetCurrentAdminInfoRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{2}
}

// 获取当前管理员响应
type AdminServiceGetCurrentAdminInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Admin         *AdminInfo             `protobuf:"bytes,1,opt,name=admin,proto3" json:"admin,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminServiceGetCurrentAdminInfoResponse) Reset() {
	*x = AdminServiceGetCurrentAdminInfoResponse{}
	mi := &file_proto_v1_admin_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminServiceGetCurrentAdminInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminServiceGetCurrentAdminInfoResponse) ProtoMessage() {}

func (x *AdminServiceGetCurrentAdminInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminServiceGetCurrentAdminInfoResponse.ProtoReflect.Descriptor instead.
func (*AdminServiceGetCurrentAdminInfoResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{3}
}

func (x *AdminServiceGetCurrentAdminInfoResponse) GetAdmin() *AdminInfo {
	if x != nil {
		return x.Admin
	}
	return nil
}

// 系统初始化请求
type AdminServiceInitializeSystemRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"` // 管理员用户名
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"` // 管理员密码
	Nickname      string                 `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"` // 管理员昵称
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminServiceInitializeSystemRequest) Reset() {
	*x = AdminServiceInitializeSystemRequest{}
	mi := &file_proto_v1_admin_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminServiceInitializeSystemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminServiceInitializeSystemRequest) ProtoMessage() {}

func (x *AdminServiceInitializeSystemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminServiceInitializeSystemRequest.ProtoReflect.Descriptor instead.
func (*AdminServiceInitializeSystemRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{4}
}

func (x *AdminServiceInitializeSystemRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *AdminServiceInitializeSystemRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *AdminServiceInitializeSystemRequest) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

// 系统初始化响应
type AdminServiceInitializeSystemResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Admin         *AdminInfo             `protobuf:"bytes,1,opt,name=admin,proto3" json:"admin,omitempty"`                                   // 创建的管理员信息
	AccessToken   string                 `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`    // 访问令牌
	RefreshToken  string                 `protobuf:"bytes,3,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"` // 刷新令牌
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminServiceInitializeSystemResponse) Reset() {
	*x = AdminServiceInitializeSystemResponse{}
	mi := &file_proto_v1_admin_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminServiceInitializeSystemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminServiceInitializeSystemResponse) ProtoMessage() {}

func (x *AdminServiceInitializeSystemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminServiceInitializeSystemResponse.ProtoReflect.Descriptor instead.
func (*AdminServiceInitializeSystemResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{5}
}

func (x *AdminServiceInitializeSystemResponse) GetAdmin() *AdminInfo {
	if x != nil {
		return x.Admin
	}
	return nil
}

func (x *AdminServiceInitializeSystemResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *AdminServiceInitializeSystemResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

// 管理员登录请求
type AdminServiceAdminLoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Username      string                 `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminServiceAdminLoginRequest) Reset() {
	*x = AdminServiceAdminLoginRequest{}
	mi := &file_proto_v1_admin_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminServiceAdminLoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminServiceAdminLoginRequest) ProtoMessage() {}

func (x *AdminServiceAdminLoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminServiceAdminLoginRequest.ProtoReflect.Descriptor instead.
func (*AdminServiceAdminLoginRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{6}
}

func (x *AdminServiceAdminLoginRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *AdminServiceAdminLoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// 管理员登录响应
type AdminServiceAdminLoginResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccessToken   string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken  string                 `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	Admin         *AdminInfo             `protobuf:"bytes,3,opt,name=admin,proto3" json:"admin,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminServiceAdminLoginResponse) Reset() {
	*x = AdminServiceAdminLoginResponse{}
	mi := &file_proto_v1_admin_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminServiceAdminLoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminServiceAdminLoginResponse) ProtoMessage() {}

func (x *AdminServiceAdminLoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminServiceAdminLoginResponse.ProtoReflect.Descriptor instead.
func (*AdminServiceAdminLoginResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{7}
}

func (x *AdminServiceAdminLoginResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *AdminServiceAdminLoginResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *AdminServiceAdminLoginResponse) GetAdmin() *AdminInfo {
	if x != nil {
		return x.Admin
	}
	return nil
}

// 管理员信息
type AdminInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Nickname      string                 `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Roles         []string               `protobuf:"bytes,4,rep,name=roles,proto3" json:"roles,omitempty"`
	CreatedAt     int64                  `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     int64                  `protobuf:"varint,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminInfo) Reset() {
	*x = AdminInfo{}
	mi := &file_proto_v1_admin_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminInfo) ProtoMessage() {}

func (x *AdminInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminInfo.ProtoReflect.Descriptor instead.
func (*AdminInfo) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{8}
}

func (x *AdminInfo) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *AdminInfo) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *AdminInfo) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *AdminInfo) GetRoles() []string {
	if x != nil {
		return x.Roles
	}
	return nil
}

func (x *AdminInfo) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *AdminInfo) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

// 刷新令牌请求
type AdminServiceRefreshTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RefreshToken  string                 `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminServiceRefreshTokenRequest) Reset() {
	*x = AdminServiceRefreshTokenRequest{}
	mi := &file_proto_v1_admin_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminServiceRefreshTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminServiceRefreshTokenRequest) ProtoMessage() {}

func (x *AdminServiceRefreshTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminServiceRefreshTokenRequest.ProtoReflect.Descriptor instead.
func (*AdminServiceRefreshTokenRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{9}
}

func (x *AdminServiceRefreshTokenRequest) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

// 刷新令牌响应
type AdminServiceRefreshTokenResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AccessToken   string                 `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken  string                 `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AdminServiceRefreshTokenResponse) Reset() {
	*x = AdminServiceRefreshTokenResponse{}
	mi := &file_proto_v1_admin_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AdminServiceRefreshTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AdminServiceRefreshTokenResponse) ProtoMessage() {}

func (x *AdminServiceRefreshTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_admin_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AdminServiceRefreshTokenResponse.ProtoReflect.Descriptor instead.
func (*AdminServiceRefreshTokenResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_admin_proto_rawDescGZIP(), []int{10}
}

func (x *AdminServiceRefreshTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *AdminServiceRefreshTokenResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

var File_proto_v1_admin_proto protoreflect.FileDescriptor

var file_proto_v1_admin_proto_rawDesc = string([]byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x49,
	0x0a, 0x25, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x69, 0x6e, 0x69, 0x74, 0x69,
	0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x69, 0x6e,
	0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x64, 0x22, 0x26, 0x0a, 0x24, 0x41, 0x64, 0x6d,
	0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53, 0x79,
	0x73, 0x74, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x28, 0x0a, 0x26, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x47, 0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x54, 0x0a, 0x27, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x47, 0x65, 0x74, 0x43, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x61, 0x64, 0x6d, 0x69,
	0x6e, 0x22, 0x79, 0x0a, 0x23, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x53, 0x79, 0x73, 0x74, 0x65,
	0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x99, 0x01, 0x0a,
	0x24, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e, 0x69,
	0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x57, 0x0a, 0x1d, 0x41, 0x64, 0x6d, 0x69,
	0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65,
	0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72,
	0x64, 0x22, 0x93, 0x01, 0x0a, 0x1e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65,
	0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x29, 0x0a, 0x05,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x05, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x22, 0xa7, 0x01, 0x0a, 0x09, 0x41, 0x64, 0x6d, 0x69,
	0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x72, 0x6f,
	0x6c, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x22, 0x46, 0x0a, 0x1f, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x66,
	0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x6a, 0x0a, 0x20, 0x41, 0x64, 0x6d,
	0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a,
	0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x32, 0xf2, 0x05, 0x0a, 0x0c, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x99, 0x01, 0x0a, 0x11, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x23, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x1d, 0x12, 0x1b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61,
	0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x9d, 0x01, 0x0a, 0x10, 0x49, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a,
	0x65, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x2d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49,
	0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76,
	0x31, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x6e,
	0x69, 0x74, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x24, 0x3a, 0x01,
	0x2a, 0x22, 0x1f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x2f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x2f, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x69,
	0x7a, 0x65, 0x12, 0x7f, 0x0a, 0x0a, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x12, 0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x6d, 0x69,
	0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x3a, 0x01, 0x2a, 0x22, 0x13,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x6c, 0x6f,
	0x67, 0x69, 0x6e, 0x12, 0x8d, 0x01, 0x0a, 0x0c, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x29, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e,
	0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x66, 0x72,
	0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x26, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x20, 0x3a, 0x01, 0x2a, 0x22, 0x1b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x2d, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x94, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x30, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x47, 0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x41, 0x64, 0x6d,
	0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x47, 0x65, 0x74, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x41,
	0x64, 0x6d, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x18, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x12, 0x12, 0x10, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76,
	0x31, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x6d, 0x65, 0x42, 0x31, 0x5a, 0x28, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x61, 0x7a, 0x79, 0x6a, 0x65, 0x61,
	0x6e, 0x2f, 0x73, 0x6c, 0x61, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x76, 0x31, 0x3b, 0x70, 0x62, 0xba, 0x02, 0x04, 0x53, 0x4c, 0x41, 0x32, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_v1_admin_proto_rawDescOnce sync.Once
	file_proto_v1_admin_proto_rawDescData []byte
)

func file_proto_v1_admin_proto_rawDescGZIP() []byte {
	file_proto_v1_admin_proto_rawDescOnce.Do(func() {
		file_proto_v1_admin_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_v1_admin_proto_rawDesc), len(file_proto_v1_admin_proto_rawDesc)))
	})
	return file_proto_v1_admin_proto_rawDescData
}

var file_proto_v1_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_proto_v1_admin_proto_goTypes = []any{
	(*AdminServiceCheckSystemStatusResponse)(nil),   // 0: proto.v1.AdminServiceCheckSystemStatusResponse
	(*AdminServiceCheckSystemStatusRequest)(nil),    // 1: proto.v1.AdminServiceCheckSystemStatusRequest
	(*AdminServiceGetCurrentAdminInfoRequest)(nil),  // 2: proto.v1.AdminServiceGetCurrentAdminInfoRequest
	(*AdminServiceGetCurrentAdminInfoResponse)(nil), // 3: proto.v1.AdminServiceGetCurrentAdminInfoResponse
	(*AdminServiceInitializeSystemRequest)(nil),     // 4: proto.v1.AdminServiceInitializeSystemRequest
	(*AdminServiceInitializeSystemResponse)(nil),    // 5: proto.v1.AdminServiceInitializeSystemResponse
	(*AdminServiceAdminLoginRequest)(nil),           // 6: proto.v1.AdminServiceAdminLoginRequest
	(*AdminServiceAdminLoginResponse)(nil),          // 7: proto.v1.AdminServiceAdminLoginResponse
	(*AdminInfo)(nil),                               // 8: proto.v1.AdminInfo
	(*AdminServiceRefreshTokenRequest)(nil),         // 9: proto.v1.AdminServiceRefreshTokenRequest
	(*AdminServiceRefreshTokenResponse)(nil),        // 10: proto.v1.AdminServiceRefreshTokenResponse
}
var file_proto_v1_admin_proto_depIdxs = []int32{
	8,  // 0: proto.v1.AdminServiceGetCurrentAdminInfoResponse.admin:type_name -> proto.v1.AdminInfo
	8,  // 1: proto.v1.AdminServiceInitializeSystemResponse.admin:type_name -> proto.v1.AdminInfo
	8,  // 2: proto.v1.AdminServiceAdminLoginResponse.admin:type_name -> proto.v1.AdminInfo
	1,  // 3: proto.v1.AdminService.CheckSystemStatus:input_type -> proto.v1.AdminServiceCheckSystemStatusRequest
	4,  // 4: proto.v1.AdminService.InitializeSystem:input_type -> proto.v1.AdminServiceInitializeSystemRequest
	6,  // 5: proto.v1.AdminService.AdminLogin:input_type -> proto.v1.AdminServiceAdminLoginRequest
	9,  // 6: proto.v1.AdminService.RefreshToken:input_type -> proto.v1.AdminServiceRefreshTokenRequest
	2,  // 7: proto.v1.AdminService.GetCurrentAdminInfo:input_type -> proto.v1.AdminServiceGetCurrentAdminInfoRequest
	0,  // 8: proto.v1.AdminService.CheckSystemStatus:output_type -> proto.v1.AdminServiceCheckSystemStatusResponse
	5,  // 9: proto.v1.AdminService.InitializeSystem:output_type -> proto.v1.AdminServiceInitializeSystemResponse
	7,  // 10: proto.v1.AdminService.AdminLogin:output_type -> proto.v1.AdminServiceAdminLoginResponse
	10, // 11: proto.v1.AdminService.RefreshToken:output_type -> proto.v1.AdminServiceRefreshTokenResponse
	3,  // 12: proto.v1.AdminService.GetCurrentAdminInfo:output_type -> proto.v1.AdminServiceGetCurrentAdminInfoResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_proto_v1_admin_proto_init() }
func file_proto_v1_admin_proto_init() {
	if File_proto_v1_admin_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_v1_admin_proto_rawDesc), len(file_proto_v1_admin_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_v1_admin_proto_goTypes,
		DependencyIndexes: file_proto_v1_admin_proto_depIdxs,
		MessageInfos:      file_proto_v1_admin_proto_msgTypes,
	}.Build()
	File_proto_v1_admin_proto = out.File
	file_proto_v1_admin_proto_goTypes = nil
	file_proto_v1_admin_proto_depIdxs = nil
}

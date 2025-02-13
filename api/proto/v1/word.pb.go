// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.21.12
// source: proto/v1/word.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type Word struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Spelling      string                 `protobuf:"bytes,2,opt,name=spelling,proto3" json:"spelling,omitempty"`
	Pronunciation string                 `protobuf:"bytes,3,opt,name=pronunciation,proto3" json:"pronunciation,omitempty"`
	Definitions   []string               `protobuf:"bytes,4,rep,name=definitions,proto3" json:"definitions,omitempty"`
	Examples      []string               `protobuf:"bytes,5,rep,name=examples,proto3" json:"examples,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Word) Reset() {
	*x = Word{}
	mi := &file_proto_v1_word_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Word) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Word) ProtoMessage() {}

func (x *Word) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_word_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Word.ProtoReflect.Descriptor instead.
func (*Word) Descriptor() ([]byte, []int) {
	return file_proto_v1_word_proto_rawDescGZIP(), []int{0}
}

func (x *Word) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Word) GetSpelling() string {
	if x != nil {
		return x.Spelling
	}
	return ""
}

func (x *Word) GetPronunciation() string {
	if x != nil {
		return x.Pronunciation
	}
	return ""
}

func (x *Word) GetDefinitions() []string {
	if x != nil {
		return x.Definitions
	}
	return nil
}

func (x *Word) GetExamples() []string {
	if x != nil {
		return x.Examples
	}
	return nil
}

func (x *Word) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Word) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type GetWordRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	WordId        uint32                 `protobuf:"varint,1,opt,name=word_id,json=wordId,proto3" json:"word_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetWordRequest) Reset() {
	*x = GetWordRequest{}
	mi := &file_proto_v1_word_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetWordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWordRequest) ProtoMessage() {}

func (x *GetWordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_word_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWordRequest.ProtoReflect.Descriptor instead.
func (*GetWordRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_word_proto_rawDescGZIP(), []int{1}
}

func (x *GetWordRequest) GetWordId() uint32 {
	if x != nil {
		return x.WordId
	}
	return 0
}

type GetWordResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Word          *Word                  `protobuf:"bytes,1,opt,name=word,proto3" json:"word,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetWordResponse) Reset() {
	*x = GetWordResponse{}
	mi := &file_proto_v1_word_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetWordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetWordResponse) ProtoMessage() {}

func (x *GetWordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_word_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetWordResponse.ProtoReflect.Descriptor instead.
func (*GetWordResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_word_proto_rawDescGZIP(), []int{2}
}

func (x *GetWordResponse) GetWord() *Word {
	if x != nil {
		return x.Word
	}
	return nil
}

type ListWordsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          uint32                 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      uint32                 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListWordsRequest) Reset() {
	*x = ListWordsRequest{}
	mi := &file_proto_v1_word_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListWordsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWordsRequest) ProtoMessage() {}

func (x *ListWordsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_word_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWordsRequest.ProtoReflect.Descriptor instead.
func (*ListWordsRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_word_proto_rawDescGZIP(), []int{3}
}

func (x *ListWordsRequest) GetPage() uint32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListWordsRequest) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type ListWordsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Words         []*Word                `protobuf:"bytes,1,rep,name=words,proto3" json:"words,omitempty"`
	Total         uint32                 `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListWordsResponse) Reset() {
	*x = ListWordsResponse{}
	mi := &file_proto_v1_word_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListWordsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListWordsResponse) ProtoMessage() {}

func (x *ListWordsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_word_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListWordsResponse.ProtoReflect.Descriptor instead.
func (*ListWordsResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_word_proto_rawDescGZIP(), []int{4}
}

func (x *ListWordsResponse) GetWords() []*Word {
	if x != nil {
		return x.Words
	}
	return nil
}

func (x *ListWordsResponse) GetTotal() uint32 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_proto_v1_word_proto protoreflect.FileDescriptor

var file_proto_v1_word_proto_rawDesc = string([]byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x6f, 0x72, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x6c, 0x61, 0x32, 0x2e, 0x76, 0x31, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x8c, 0x02, 0x0a, 0x04, 0x57, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x70, 0x65, 0x6c,
	0x6c, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x70, 0x65, 0x6c,
	0x6c, 0x69, 0x6e, 0x67, 0x12, 0x24, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x6e, 0x75, 0x6e, 0x63, 0x69,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x70, 0x72, 0x6f,
	0x6e, 0x75, 0x6e, 0x63, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1a, 0x0a, 0x08,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x29,
	0x0a, 0x0e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x06, 0x77, 0x6f, 0x72, 0x64, 0x49, 0x64, 0x22, 0x34, 0x0a, 0x0f, 0x47, 0x65, 0x74,
	0x57, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x04,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x73, 0x6c, 0x61,
	0x32, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x6f, 0x72, 0x64, 0x52, 0x04, 0x77, 0x6f, 0x72, 0x64, 0x22,
	0x43, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65,
	0x53, 0x69, 0x7a, 0x65, 0x22, 0x4e, 0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x64,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x05, 0x77, 0x6f, 0x72,
	0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x73, 0x6c, 0x61, 0x32, 0x2e,
	0x76, 0x31, 0x2e, 0x57, 0x6f, 0x72, 0x64, 0x52, 0x05, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x32, 0x93, 0x01, 0x0a, 0x0b, 0x57, 0x6f, 0x72, 0x64, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x3e, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x12,
	0x17, 0x2e, 0x73, 0x6c, 0x61, 0x32, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72,
	0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x6c, 0x61, 0x32, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x64,
	0x73, 0x12, 0x19, 0x2e, 0x73, 0x6c, 0x61, 0x32, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x57, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x73,
	0x6c, 0x61, 0x32, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x57, 0x6f, 0x72, 0x64, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6c, 0x69, 0x75, 0x7a, 0x68, 0x65, 0x6e,
	0x32, 0x31, 0x2f, 0x73, 0x6c, 0x61, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x76, 0x31, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_proto_v1_word_proto_rawDescOnce sync.Once
	file_proto_v1_word_proto_rawDescData []byte
)

func file_proto_v1_word_proto_rawDescGZIP() []byte {
	file_proto_v1_word_proto_rawDescOnce.Do(func() {
		file_proto_v1_word_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_v1_word_proto_rawDesc), len(file_proto_v1_word_proto_rawDesc)))
	})
	return file_proto_v1_word_proto_rawDescData
}

var file_proto_v1_word_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_v1_word_proto_goTypes = []any{
	(*Word)(nil),                  // 0: sla2.v1.Word
	(*GetWordRequest)(nil),        // 1: sla2.v1.GetWordRequest
	(*GetWordResponse)(nil),       // 2: sla2.v1.GetWordResponse
	(*ListWordsRequest)(nil),      // 3: sla2.v1.ListWordsRequest
	(*ListWordsResponse)(nil),     // 4: sla2.v1.ListWordsResponse
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_proto_v1_word_proto_depIdxs = []int32{
	5, // 0: sla2.v1.Word.created_at:type_name -> google.protobuf.Timestamp
	5, // 1: sla2.v1.Word.updated_at:type_name -> google.protobuf.Timestamp
	0, // 2: sla2.v1.GetWordResponse.word:type_name -> sla2.v1.Word
	0, // 3: sla2.v1.ListWordsResponse.words:type_name -> sla2.v1.Word
	1, // 4: sla2.v1.WordService.GetWord:input_type -> sla2.v1.GetWordRequest
	3, // 5: sla2.v1.WordService.ListWords:input_type -> sla2.v1.ListWordsRequest
	2, // 6: sla2.v1.WordService.GetWord:output_type -> sla2.v1.GetWordResponse
	4, // 7: sla2.v1.WordService.ListWords:output_type -> sla2.v1.ListWordsResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_v1_word_proto_init() }
func file_proto_v1_word_proto_init() {
	if File_proto_v1_word_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_v1_word_proto_rawDesc), len(file_proto_v1_word_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_v1_word_proto_goTypes,
		DependencyIndexes: file_proto_v1_word_proto_depIdxs,
		MessageInfos:      file_proto_v1_word_proto_msgTypes,
	}.Build()
	File_proto_v1_word_proto = out.File
	file_proto_v1_word_proto_goTypes = nil
	file_proto_v1_word_proto_depIdxs = nil
}

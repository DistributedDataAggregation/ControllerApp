// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.21.12
// source: query_response.proto

package protomodels

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type QueryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Guid   string   `protobuf:"bytes,1,opt,name=guid,proto3" json:"guid,omitempty"`
	Error  *Error   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	Values []*Value `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *QueryResponse) Reset() {
	*x = QueryResponse{}
	mi := &file_query_response_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryResponse) ProtoMessage() {}

func (x *QueryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_query_response_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryResponse.ProtoReflect.Descriptor instead.
func (*QueryResponse) Descriptor() ([]byte, []int) {
	return file_query_response_proto_rawDescGZIP(), []int{0}
}

func (x *QueryResponse) GetGuid() string {
	if x != nil {
		return x.Guid
	}
	return ""
}

func (x *QueryResponse) GetError() *Error {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *QueryResponse) GetValues() []*Value {
	if x != nil {
		return x.Values
	}
	return nil
}

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupingValue string           `protobuf:"bytes,1,opt,name=grouping_value,json=groupingValue,proto3" json:"grouping_value,omitempty"`
	Results       []*PartialResult `protobuf:"bytes,2,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *Value) Reset() {
	*x = Value{}
	mi := &file_query_response_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_query_response_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_query_response_proto_rawDescGZIP(), []int{1}
}

func (x *Value) GetGroupingValue() string {
	if x != nil {
		return x.GroupingValue
	}
	return ""
}

func (x *Value) GetResults() []*PartialResult {
	if x != nil {
		return x.Results
	}
	return nil
}

type PartialResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsNull bool  `protobuf:"varint,1,opt,name=is_null,json=isNull,proto3" json:"is_null,omitempty"`
	Value  int64 `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	Count  int64 `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *PartialResult) Reset() {
	*x = PartialResult{}
	mi := &file_query_response_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PartialResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartialResult) ProtoMessage() {}

func (x *PartialResult) ProtoReflect() protoreflect.Message {
	mi := &file_query_response_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartialResult.ProtoReflect.Descriptor instead.
func (*PartialResult) Descriptor() ([]byte, []int) {
	return file_query_response_proto_rawDescGZIP(), []int{2}
}

func (x *PartialResult) GetIsNull() bool {
	if x != nil {
		return x.IsNull
	}
	return false
}

func (x *PartialResult) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *PartialResult) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message      string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	InnerMessage string `protobuf:"bytes,2,opt,name=inner_message,json=innerMessage,proto3" json:"inner_message,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	mi := &file_query_response_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_query_response_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_query_response_proto_rawDescGZIP(), []int{3}
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Error) GetInnerMessage() string {
	if x != nil {
		return x.InnerMessage
	}
	return ""
}

var File_query_response_proto protoreflect.FileDescriptor

var file_query_response_proto_rawDesc = []byte{
	0x0a, 0x14, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x61, 0x0a, 0x0d, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x67, 0x75, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1e, 0x0a, 0x06, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x58, 0x0a, 0x05, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x69, 0x6e, 0x67, 0x5f, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x50, 0x61, 0x72,
	0x74, 0x69, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x73, 0x22, 0x54, 0x0a, 0x0d, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x6e, 0x75, 0x6c, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x4e, 0x75, 0x6c, 0x6c, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x46, 0x0a, 0x05, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x23, 0x0a, 0x0d,
	0x69, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x42, 0x0e, 0x5a, 0x0c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_query_response_proto_rawDescOnce sync.Once
	file_query_response_proto_rawDescData = file_query_response_proto_rawDesc
)

func file_query_response_proto_rawDescGZIP() []byte {
	file_query_response_proto_rawDescOnce.Do(func() {
		file_query_response_proto_rawDescData = protoimpl.X.CompressGZIP(file_query_response_proto_rawDescData)
	})
	return file_query_response_proto_rawDescData
}

var file_query_response_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_query_response_proto_goTypes = []any{
	(*QueryResponse)(nil), // 0: QueryResponse
	(*Value)(nil),         // 1: Value
	(*PartialResult)(nil), // 2: PartialResult
	(*Error)(nil),         // 3: Error
}
var file_query_response_proto_depIdxs = []int32{
	3, // 0: QueryResponse.error:type_name -> Error
	1, // 1: QueryResponse.values:type_name -> Value
	2, // 2: Value.results:type_name -> PartialResult
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_query_response_proto_init() }
func file_query_response_proto_init() {
	if File_query_response_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_query_response_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_query_response_proto_goTypes,
		DependencyIndexes: file_query_response_proto_depIdxs,
		MessageInfos:      file_query_response_proto_msgTypes,
	}.Build()
	File_query_response_proto = out.File
	file_query_response_proto_rawDesc = nil
	file_query_response_proto_goTypes = nil
	file_query_response_proto_depIdxs = nil
}

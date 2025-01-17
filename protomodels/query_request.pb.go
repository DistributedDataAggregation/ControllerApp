// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.21.12
// source: query_request.proto

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

type QueryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Guid         string               `protobuf:"bytes,1,opt,name=guid,proto3" json:"guid,omitempty"`
	FilesNames   []string             `protobuf:"bytes,2,rep,name=files_names,json=filesNames,proto3" json:"files_names,omitempty"`
	GroupColumns []string             `protobuf:"bytes,3,rep,name=group_columns,json=groupColumns,proto3" json:"group_columns,omitempty"`
	Select       []*Select            `protobuf:"bytes,4,rep,name=select,proto3" json:"select,omitempty"`
	Executor     *ExecutorInformation `protobuf:"bytes,5,opt,name=executor,proto3" json:"executor,omitempty"`
}

func (x *QueryRequest) Reset() {
	*x = QueryRequest{}
	mi := &file_query_request_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRequest) ProtoMessage() {}

func (x *QueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_query_request_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRequest.ProtoReflect.Descriptor instead.
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return file_query_request_proto_rawDescGZIP(), []int{0}
}

func (x *QueryRequest) GetGuid() string {
	if x != nil {
		return x.Guid
	}
	return ""
}

func (x *QueryRequest) GetFilesNames() []string {
	if x != nil {
		return x.FilesNames
	}
	return nil
}

func (x *QueryRequest) GetGroupColumns() []string {
	if x != nil {
		return x.GroupColumns
	}
	return nil
}

func (x *QueryRequest) GetSelect() []*Select {
	if x != nil {
		return x.Select
	}
	return nil
}

func (x *QueryRequest) GetExecutor() *ExecutorInformation {
	if x != nil {
		return x.Executor
	}
	return nil
}

type Select struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Column   string    `protobuf:"bytes,1,opt,name=column,proto3" json:"column,omitempty"`
	Function Aggregate `protobuf:"varint,2,opt,name=function,proto3,enum=Aggregate" json:"function,omitempty"`
}

func (x *Select) Reset() {
	*x = Select{}
	mi := &file_query_request_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Select) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Select) ProtoMessage() {}

func (x *Select) ProtoReflect() protoreflect.Message {
	mi := &file_query_request_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Select.ProtoReflect.Descriptor instead.
func (*Select) Descriptor() ([]byte, []int) {
	return file_query_request_proto_rawDescGZIP(), []int{1}
}

func (x *Select) GetColumn() string {
	if x != nil {
		return x.Column
	}
	return ""
}

func (x *Select) GetFunction() Aggregate {
	if x != nil {
		return x.Function
	}
	return Aggregate_Minimum
}

type ExecutorInformation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsCurrentNodeMain bool   `protobuf:"varint,1,opt,name=is_current_node_main,json=isCurrentNodeMain,proto3" json:"is_current_node_main,omitempty"`
	MainIpAddress     string `protobuf:"bytes,2,opt,name=main_ip_address,json=mainIpAddress,proto3" json:"main_ip_address,omitempty"`
	MainPort          int32  `protobuf:"varint,3,opt,name=main_port,json=mainPort,proto3" json:"main_port,omitempty"`
	ExecutorsCount    int32  `protobuf:"varint,4,opt,name=executors_count,json=executorsCount,proto3" json:"executors_count,omitempty"`
}

func (x *ExecutorInformation) Reset() {
	*x = ExecutorInformation{}
	mi := &file_query_request_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExecutorInformation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecutorInformation) ProtoMessage() {}

func (x *ExecutorInformation) ProtoReflect() protoreflect.Message {
	mi := &file_query_request_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecutorInformation.ProtoReflect.Descriptor instead.
func (*ExecutorInformation) Descriptor() ([]byte, []int) {
	return file_query_request_proto_rawDescGZIP(), []int{2}
}

func (x *ExecutorInformation) GetIsCurrentNodeMain() bool {
	if x != nil {
		return x.IsCurrentNodeMain
	}
	return false
}

func (x *ExecutorInformation) GetMainIpAddress() string {
	if x != nil {
		return x.MainIpAddress
	}
	return ""
}

func (x *ExecutorInformation) GetMainPort() int32 {
	if x != nil {
		return x.MainPort
	}
	return 0
}

func (x *ExecutorInformation) GetExecutorsCount() int32 {
	if x != nil {
		return x.ExecutorsCount
	}
	return 0
}

var File_query_request_proto protoreflect.FileDescriptor

var file_query_request_proto_rawDesc = []byte{
	0x0a, 0x13, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbb, 0x01, 0x0a, 0x0c, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x67, 0x75, 0x69, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x66,
	0x69, 0x6c, 0x65, 0x73, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x23, 0x0a, 0x0d,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0c, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e,
	0x73, 0x12, 0x1f, 0x0a, 0x06, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x07, 0x2e, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x52, 0x06, 0x73, 0x65, 0x6c, 0x65,
	0x63, 0x74, 0x12, 0x30, 0x0a, 0x08, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x65, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x6f, 0x72, 0x22, 0x48, 0x0a, 0x06, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x12, 0x26, 0x0a, 0x08, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x52, 0x08, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xb4,
	0x01, 0x0a, 0x13, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x72,
	0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x14, 0x69, 0x73, 0x5f, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x11, 0x69, 0x73, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x4e,
	0x6f, 0x64, 0x65, 0x4d, 0x61, 0x69, 0x6e, 0x12, 0x26, 0x0a, 0x0f, 0x6d, 0x61, 0x69, 0x6e, 0x5f,
	0x69, 0x70, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x6d, 0x61, 0x69, 0x6e, 0x49, 0x70, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12,
	0x1b, 0x0a, 0x09, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x6d, 0x61, 0x69, 0x6e, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x27, 0x0a, 0x0f,
	0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x6f, 0x72, 0x73,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x0e, 0x5a, 0x0c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_query_request_proto_rawDescOnce sync.Once
	file_query_request_proto_rawDescData = file_query_request_proto_rawDesc
)

func file_query_request_proto_rawDescGZIP() []byte {
	file_query_request_proto_rawDescOnce.Do(func() {
		file_query_request_proto_rawDescData = protoimpl.X.CompressGZIP(file_query_request_proto_rawDescData)
	})
	return file_query_request_proto_rawDescData
}

var file_query_request_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_query_request_proto_goTypes = []any{
	(*QueryRequest)(nil),        // 0: QueryRequest
	(*Select)(nil),              // 1: Select
	(*ExecutorInformation)(nil), // 2: ExecutorInformation
	(Aggregate)(0),              // 3: Aggregate
}
var file_query_request_proto_depIdxs = []int32{
	1, // 0: QueryRequest.select:type_name -> Select
	2, // 1: QueryRequest.executor:type_name -> ExecutorInformation
	3, // 2: Select.function:type_name -> Aggregate
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_query_request_proto_init() }
func file_query_request_proto_init() {
	if File_query_request_proto != nil {
		return
	}
	file_aggregate_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_query_request_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_query_request_proto_goTypes,
		DependencyIndexes: file_query_request_proto_depIdxs,
		MessageInfos:      file_query_request_proto_msgTypes,
	}.Build()
	File_query_request_proto = out.File
	file_query_request_proto_rawDesc = nil
	file_query_request_proto_goTypes = nil
	file_query_request_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
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

type ResultType int32

const (
	ResultType_UNKNOWN ResultType = 0
	ResultType_INT     ResultType = 1
	ResultType_FLOAT   ResultType = 2
	ResultType_DOUBLE  ResultType = 3
)

// Enum value maps for ResultType.
var (
	ResultType_name = map[int32]string{
		0: "UNKNOWN",
		1: "INT",
		2: "FLOAT",
		3: "DOUBLE",
	}
	ResultType_value = map[string]int32{
		"UNKNOWN": 0,
		"INT":     1,
		"FLOAT":   2,
		"DOUBLE":  3,
	}
)

func (x ResultType) Enum() *ResultType {
	p := new(ResultType)
	*p = x
	return p
}

func (x ResultType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResultType) Descriptor() protoreflect.EnumDescriptor {
	return file_query_response_proto_enumTypes[0].Descriptor()
}

func (ResultType) Type() protoreflect.EnumType {
	return &file_query_response_proto_enumTypes[0]
}

func (x ResultType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResultType.Descriptor instead.
func (ResultType) EnumDescriptor() ([]byte, []int) {
	return file_query_response_proto_rawDescGZIP(), []int{0}
}

type QueryResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Guid          string                 `protobuf:"bytes,1,opt,name=guid,proto3" json:"guid,omitempty"`
	Error         *Error                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	Values        []*Value               `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
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
	state         protoimpl.MessageState `protogen:"open.v1"`
	GroupingValue string                 `protobuf:"bytes,1,opt,name=grouping_value,json=groupingValue,proto3" json:"grouping_value,omitempty"`
	Results       []*PartialResult       `protobuf:"bytes,2,rep,name=results,proto3" json:"results,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
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
	state  protoimpl.MessageState `protogen:"open.v1"`
	IsNull bool                   `protobuf:"varint,1,opt,name=is_null,json=isNull,proto3" json:"is_null,omitempty"`
	Type   ResultType             `protobuf:"varint,2,opt,name=type,proto3,enum=ResultType" json:"type,omitempty"`
	// Types that are valid to be assigned to Value:
	//
	//	*PartialResult_IntValue
	//	*PartialResult_FloatValue
	//	*PartialResult_DoubleValue
	Value         isPartialResult_Value `protobuf_oneof:"value"`
	Count         int64                 `protobuf:"varint,7,opt,name=count,proto3" json:"count,omitempty"`
	Function      Aggregate             `protobuf:"varint,8,opt,name=function,proto3,enum=Aggregate" json:"function,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
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

func (x *PartialResult) GetType() ResultType {
	if x != nil {
		return x.Type
	}
	return ResultType_UNKNOWN
}

func (x *PartialResult) GetValue() isPartialResult_Value {
	if x != nil {
		return x.Value
	}
	return nil
}

func (x *PartialResult) GetIntValue() int64 {
	if x != nil {
		if x, ok := x.Value.(*PartialResult_IntValue); ok {
			return x.IntValue
		}
	}
	return 0
}

func (x *PartialResult) GetFloatValue() float32 {
	if x != nil {
		if x, ok := x.Value.(*PartialResult_FloatValue); ok {
			return x.FloatValue
		}
	}
	return 0
}

func (x *PartialResult) GetDoubleValue() float64 {
	if x != nil {
		if x, ok := x.Value.(*PartialResult_DoubleValue); ok {
			return x.DoubleValue
		}
	}
	return 0
}

func (x *PartialResult) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *PartialResult) GetFunction() Aggregate {
	if x != nil {
		return x.Function
	}
	return Aggregate_Minimum
}

type isPartialResult_Value interface {
	isPartialResult_Value()
}

type PartialResult_IntValue struct {
	IntValue int64 `protobuf:"varint,3,opt,name=int_value,json=intValue,proto3,oneof"`
}

type PartialResult_FloatValue struct {
	FloatValue float32 `protobuf:"fixed32,4,opt,name=float_value,json=floatValue,proto3,oneof"`
}

type PartialResult_DoubleValue struct {
	DoubleValue float64 `protobuf:"fixed64,5,opt,name=double_value,json=doubleValue,proto3,oneof"`
}

func (*PartialResult_IntValue) isPartialResult_Value() {}

func (*PartialResult_FloatValue) isPartialResult_Value() {}

func (*PartialResult_DoubleValue) isPartialResult_Value() {}

type Error struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	InnerMessage  string                 `protobuf:"bytes,2,opt,name=inner_message,json=innerMessage,proto3" json:"inner_message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
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
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x61, 0x0a, 0x0d, 0x51, 0x75, 0x65, 0x72, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x67, 0x75, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x1e, 0x0a, 0x06, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x58, 0x0a, 0x05, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x69, 0x6e, 0x67, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x28, 0x0a, 0x07, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x50, 0x61,
	0x72, 0x74, 0x69, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x22, 0xf7, 0x01, 0x0a, 0x0d, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x6e, 0x75, 0x6c,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x4e, 0x75, 0x6c, 0x6c, 0x12,
	0x1f, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x1d, 0x0a, 0x09, 0x69, 0x6e, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x21, 0x0a, 0x0b, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x02, 0x48, 0x00, 0x52, 0x0a, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x6f, 0x75, 0x62,
	0x6c, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x26, 0x0a,
	0x08, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0a, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x52, 0x08, 0x66, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x46,
	0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x23, 0x0a, 0x0d, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x69, 0x6e, 0x6e, 0x65, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x39, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10,
	0x00, 0x12, 0x07, 0x0a, 0x03, 0x49, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x4c,
	0x4f, 0x41, 0x54, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x4f, 0x55, 0x42, 0x4c, 0x45, 0x10,
	0x03, 0x42, 0x0e, 0x5a, 0x0c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
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

var file_query_response_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_query_response_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_query_response_proto_goTypes = []any{
	(ResultType)(0),       // 0: ResultType
	(*QueryResponse)(nil), // 1: QueryResponse
	(*Value)(nil),         // 2: Value
	(*PartialResult)(nil), // 3: PartialResult
	(*Error)(nil),         // 4: Error
	(Aggregate)(0),        // 5: Aggregate
}
var file_query_response_proto_depIdxs = []int32{
	4, // 0: QueryResponse.error:type_name -> Error
	2, // 1: QueryResponse.values:type_name -> Value
	3, // 2: Value.results:type_name -> PartialResult
	0, // 3: PartialResult.type:type_name -> ResultType
	5, // 4: PartialResult.function:type_name -> Aggregate
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_query_response_proto_init() }
func file_query_response_proto_init() {
	if File_query_response_proto != nil {
		return
	}
	file_aggregate_proto_init()
	file_query_response_proto_msgTypes[2].OneofWrappers = []any{
		(*PartialResult_IntValue)(nil),
		(*PartialResult_FloatValue)(nil),
		(*PartialResult_DoubleValue)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_query_response_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_query_response_proto_goTypes,
		DependencyIndexes: file_query_response_proto_depIdxs,
		EnumInfos:         file_query_response_proto_enumTypes,
		MessageInfos:      file_query_response_proto_msgTypes,
	}.Build()
	File_query_response_proto = out.File
	file_query_response_proto_rawDesc = nil
	file_query_response_proto_goTypes = nil
	file_query_response_proto_depIdxs = nil
}

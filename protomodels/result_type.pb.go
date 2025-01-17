// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v3.21.12
// source: result_type.proto

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
	return file_result_type_proto_enumTypes[0].Descriptor()
}

func (ResultType) Type() protoreflect.EnumType {
	return &file_result_type_proto_enumTypes[0]
}

func (x ResultType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResultType.Descriptor instead.
func (ResultType) EnumDescriptor() ([]byte, []int) {
	return file_result_type_proto_rawDescGZIP(), []int{0}
}

var File_result_type_proto protoreflect.FileDescriptor

var file_result_type_proto_rawDesc = []byte{
	0x0a, 0x11, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2a, 0x39, 0x0a, 0x0a, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x07,
	0x0a, 0x03, 0x49, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x4c, 0x4f, 0x41, 0x54,
	0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x4f, 0x55, 0x42, 0x4c, 0x45, 0x10, 0x03, 0x42, 0x0e,
	0x5a, 0x0c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_result_type_proto_rawDescOnce sync.Once
	file_result_type_proto_rawDescData = file_result_type_proto_rawDesc
)

func file_result_type_proto_rawDescGZIP() []byte {
	file_result_type_proto_rawDescOnce.Do(func() {
		file_result_type_proto_rawDescData = protoimpl.X.CompressGZIP(file_result_type_proto_rawDescData)
	})
	return file_result_type_proto_rawDescData
}

var file_result_type_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_result_type_proto_goTypes = []any{
	(ResultType)(0), // 0: ResultType
}
var file_result_type_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_result_type_proto_init() }
func file_result_type_proto_init() {
	if File_result_type_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_result_type_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_result_type_proto_goTypes,
		DependencyIndexes: file_result_type_proto_depIdxs,
		EnumInfos:         file_result_type_proto_enumTypes,
	}.Build()
	File_result_type_proto = out.File
	file_result_type_proto_rawDesc = nil
	file_result_type_proto_goTypes = nil
	file_result_type_proto_depIdxs = nil
}

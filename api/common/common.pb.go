// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: api/common/common.proto

package common

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

type Common struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 角色所在的项目id, 1~5为目前已知项目app id
	AppId uint32 `protobuf:"varint,1,opt,name=appId,proto3" json:"appId,omitempty"`
}

func (x *Common) Reset() {
	*x = Common{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_common_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Common) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Common) ProtoMessage() {}

func (x *Common) ProtoReflect() protoreflect.Message {
	mi := &file_api_common_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Common.ProtoReflect.Descriptor instead.
func (*Common) Descriptor() ([]byte, []int) {
	return file_api_common_common_proto_rawDescGZIP(), []int{0}
}

func (x *Common) GetAppId() uint32 {
	if x != nil {
		return x.AppId
	}
	return 0
}

type DownloadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName      string `protobuf:"bytes,1,opt,name=fileName,proto3" json:"fileName,omitempty"`
	ContentLength int64  `protobuf:"varint,2,opt,name=contentLength,proto3" json:"contentLength,omitempty"`
	Payload       []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *DownloadResponse) Reset() {
	*x = DownloadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_common_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadResponse) ProtoMessage() {}

func (x *DownloadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_common_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadResponse.ProtoReflect.Descriptor instead.
func (*DownloadResponse) Descriptor() ([]byte, []int) {
	return file_api_common_common_proto_rawDescGZIP(), []int{1}
}

func (x *DownloadResponse) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *DownloadResponse) GetContentLength() int64 {
	if x != nil {
		return x.ContentLength
	}
	return 0
}

func (x *DownloadResponse) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

var File_api_common_common_proto protoreflect.FileDescriptor

var file_api_common_common_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x61, 0x70, 0x69, 0x2e, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x22, 0x1e, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05,
	0x61, 0x70, 0x70, 0x49, 0x64, 0x22, 0x6e, 0x0a, 0x10, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x31, 0x32, 0x31, 0x39, 0x37, 0x39, 0x36, 0x33, 0x39, 0x35, 0x2f, 0x6d,
	0x79, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x32, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_api_common_common_proto_rawDescOnce sync.Once
	file_api_common_common_proto_rawDescData = file_api_common_common_proto_rawDesc
)

func file_api_common_common_proto_rawDescGZIP() []byte {
	file_api_common_common_proto_rawDescOnce.Do(func() {
		file_api_common_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_common_common_proto_rawDescData)
	})
	return file_api_common_common_proto_rawDescData
}

var file_api_common_common_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_common_common_proto_goTypes = []interface{}{
	(*Common)(nil),           // 0: api.common.Common
	(*DownloadResponse)(nil), // 1: api.common.DownloadResponse
}
var file_api_common_common_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_common_common_proto_init() }
func file_api_common_common_proto_init() {
	if File_api_common_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_common_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Common); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_common_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_common_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_common_common_proto_goTypes,
		DependencyIndexes: file_api_common_common_proto_depIdxs,
		MessageInfos:      file_api_common_common_proto_msgTypes,
	}.Build()
	File_api_common_common_proto = out.File
	file_api_common_common_proto_rawDesc = nil
	file_api_common_common_proto_goTypes = nil
	file_api_common_common_proto_depIdxs = nil
}

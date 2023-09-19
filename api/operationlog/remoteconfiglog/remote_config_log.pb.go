// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: api/operationlog/remoteconfiglog/remote_config_log.proto

package remoteconfiglog

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// 日志操作枚举值
type LogOperation int32

const (
	LogOperation_OPERATION_UNKNOWN        LogOperation = 0 // 占位
	LogOperation_OPERATION_CREATE         LogOperation = 1 // 创建
	LogOperation_OPERATION_UPDATE         LogOperation = 2 // 修改
	LogOperation_OPERATION_DELETE         LogOperation = 3 // 删除
	LogOperation_OPERATION_PUBLISH        LogOperation = 4 // 发布
	LogOperation_OPERATION_CANCEL_PUBLISH LogOperation = 5 // 取消发布
)

// Enum value maps for LogOperation.
var (
	LogOperation_name = map[int32]string{
		0: "OPERATION_UNKNOWN",
		1: "OPERATION_CREATE",
		2: "OPERATION_UPDATE",
		3: "OPERATION_DELETE",
		4: "OPERATION_PUBLISH",
		5: "OPERATION_CANCEL_PUBLISH",
	}
	LogOperation_value = map[string]int32{
		"OPERATION_UNKNOWN":        0,
		"OPERATION_CREATE":         1,
		"OPERATION_UPDATE":         2,
		"OPERATION_DELETE":         3,
		"OPERATION_PUBLISH":        4,
		"OPERATION_CANCEL_PUBLISH": 5,
	}
)

func (x LogOperation) Enum() *LogOperation {
	p := new(LogOperation)
	*p = x
	return p
}

func (x LogOperation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LogOperation) Descriptor() protoreflect.EnumDescriptor {
	return file_api_operationlog_remoteconfiglog_remote_config_log_proto_enumTypes[0].Descriptor()
}

func (LogOperation) Type() protoreflect.EnumType {
	return &file_api_operationlog_remoteconfiglog_remote_config_log_proto_enumTypes[0]
}

func (x LogOperation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LogOperation.Descriptor instead.
func (LogOperation) EnumDescriptor() ([]byte, []int) {
	return file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescGZIP(), []int{0}
}

type GetRemoteConfigLogListReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Appid      uint32 `protobuf:"varint,1,opt,name=appid,proto3" json:"appid,omitempty"`
	Env        string `protobuf:"bytes,2,opt,name=env,proto3" json:"env,omitempty"`
	Channel    string `protobuf:"bytes,3,opt,name=channel,proto3" json:"channel,omitempty"`
	Platform   string `protobuf:"bytes,4,opt,name=platform,proto3" json:"platform,omitempty"`
	ConfigName string `protobuf:"bytes,5,opt,name=configName,proto3" json:"configName,omitempty"`
	Operation  uint32 `protobuf:"varint,6,opt,name=operation,proto3" json:"operation,omitempty"`
	Operator   string `protobuf:"bytes,7,opt,name=operator,proto3" json:"operator,omitempty"`
	StartTime  uint64 `protobuf:"varint,8,opt,name=startTime,proto3" json:"startTime,omitempty"` // 开始时间，单位：毫秒，闭区间, [startTime, endTime)
	EndTime    uint64 `protobuf:"varint,9,opt,name=endTime,proto3" json:"endTime,omitempty"`     // 结束时间，单位：毫秒，开区间，[startTime, endTime)
}

func (x *GetRemoteConfigLogListReq) Reset() {
	*x = GetRemoteConfigLogListReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRemoteConfigLogListReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRemoteConfigLogListReq) ProtoMessage() {}

func (x *GetRemoteConfigLogListReq) ProtoReflect() protoreflect.Message {
	mi := &file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRemoteConfigLogListReq.ProtoReflect.Descriptor instead.
func (*GetRemoteConfigLogListReq) Descriptor() ([]byte, []int) {
	return file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescGZIP(), []int{0}
}

func (x *GetRemoteConfigLogListReq) GetAppid() uint32 {
	if x != nil {
		return x.Appid
	}
	return 0
}

func (x *GetRemoteConfigLogListReq) GetEnv() string {
	if x != nil {
		return x.Env
	}
	return ""
}

func (x *GetRemoteConfigLogListReq) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *GetRemoteConfigLogListReq) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *GetRemoteConfigLogListReq) GetConfigName() string {
	if x != nil {
		return x.ConfigName
	}
	return ""
}

func (x *GetRemoteConfigLogListReq) GetOperation() uint32 {
	if x != nil {
		return x.Operation
	}
	return 0
}

func (x *GetRemoteConfigLogListReq) GetOperator() string {
	if x != nil {
		return x.Operator
	}
	return ""
}

func (x *GetRemoteConfigLogListReq) GetStartTime() uint64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *GetRemoteConfigLogListReq) GetEndTime() uint64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

type GetRemoteConfigLogListRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Logs []*RemoteConfigLogDetail `protobuf:"bytes,1,rep,name=logs,proto3" json:"logs,omitempty"`
}

func (x *GetRemoteConfigLogListRsp) Reset() {
	*x = GetRemoteConfigLogListRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRemoteConfigLogListRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRemoteConfigLogListRsp) ProtoMessage() {}

func (x *GetRemoteConfigLogListRsp) ProtoReflect() protoreflect.Message {
	mi := &file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRemoteConfigLogListRsp.ProtoReflect.Descriptor instead.
func (*GetRemoteConfigLogListRsp) Descriptor() ([]byte, []int) {
	return file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescGZIP(), []int{1}
}

func (x *GetRemoteConfigLogListRsp) GetLogs() []*RemoteConfigLogDetail {
	if x != nil {
		return x.Logs
	}
	return nil
}

type RemoteConfigLogDetail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Appid        uint32 `protobuf:"varint,1,opt,name=appid,proto3" json:"appid,omitempty"`              // 游戏的业务id
	Env          string `protobuf:"bytes,2,opt,name=env,proto3" json:"env,omitempty"`                   // 环境名
	Channel      string `protobuf:"bytes,3,opt,name=channel,proto3" json:"channel,omitempty"`           // 渠道
	Platform     string `protobuf:"bytes,4,opt,name=platform,proto3" json:"platform,omitempty"`         // 平台
	ConfigName   string `protobuf:"bytes,5,opt,name=configName,proto3" json:"configName,omitempty"`     // 配置名
	Operation    uint32 `protobuf:"varint,6,opt,name=operation,proto3" json:"operation,omitempty"`      // 操作，1:创建，2:修改，3:删除，4:发布，5:取消发布
	Operator     string `protobuf:"bytes,7,opt,name=operator,proto3" json:"operator,omitempty"`         // 操作人
	UpdateBefore string `protobuf:"bytes,8,opt,name=updateBefore,proto3" json:"updateBefore,omitempty"` // 修改前的数据，json格式
	UpdateAfter  string `protobuf:"bytes,9,opt,name=updateAfter,proto3" json:"updateAfter,omitempty"`   // 修改后的数据，json格式
	CreateTime   uint64 `protobuf:"varint,10,opt,name=createTime,proto3" json:"createTime,omitempty"`   // 日志创建时间，单位：毫秒
}

func (x *RemoteConfigLogDetail) Reset() {
	*x = RemoteConfigLogDetail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoteConfigLogDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoteConfigLogDetail) ProtoMessage() {}

func (x *RemoteConfigLogDetail) ProtoReflect() protoreflect.Message {
	mi := &file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoteConfigLogDetail.ProtoReflect.Descriptor instead.
func (*RemoteConfigLogDetail) Descriptor() ([]byte, []int) {
	return file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescGZIP(), []int{2}
}

func (x *RemoteConfigLogDetail) GetAppid() uint32 {
	if x != nil {
		return x.Appid
	}
	return 0
}

func (x *RemoteConfigLogDetail) GetEnv() string {
	if x != nil {
		return x.Env
	}
	return ""
}

func (x *RemoteConfigLogDetail) GetChannel() string {
	if x != nil {
		return x.Channel
	}
	return ""
}

func (x *RemoteConfigLogDetail) GetPlatform() string {
	if x != nil {
		return x.Platform
	}
	return ""
}

func (x *RemoteConfigLogDetail) GetConfigName() string {
	if x != nil {
		return x.ConfigName
	}
	return ""
}

func (x *RemoteConfigLogDetail) GetOperation() uint32 {
	if x != nil {
		return x.Operation
	}
	return 0
}

func (x *RemoteConfigLogDetail) GetOperator() string {
	if x != nil {
		return x.Operator
	}
	return ""
}

func (x *RemoteConfigLogDetail) GetUpdateBefore() string {
	if x != nil {
		return x.UpdateBefore
	}
	return ""
}

func (x *RemoteConfigLogDetail) GetUpdateAfter() string {
	if x != nil {
		return x.UpdateAfter
	}
	return ""
}

func (x *RemoteConfigLogDetail) GetCreateTime() uint64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

var File_api_operationlog_remoteconfiglog_remote_config_log_proto protoreflect.FileDescriptor

var file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDesc = []byte{
	0x0a, 0x38, 0x61, 0x70, 0x69, 0x2f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6c,
	0x6f, 0x67, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6c,
	0x6f, 0x67, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x5f, 0x6c, 0x6f, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x20, 0x61, 0x70, 0x69, 0x2e,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6c, 0x6f, 0x67, 0x2e, 0x72, 0x65, 0x6d,
	0x6f, 0x74, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6c, 0x6f, 0x67, 0x1a, 0x1c, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8b, 0x02, 0x0a, 0x19, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4c, 0x6f,
	0x67, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x70, 0x70, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x61, 0x70, 0x70, 0x69, 0x64, 0x12, 0x10,
	0x0a, 0x03, 0x65, 0x6e, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x6e, 0x76,
	0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c,
	0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x6c,
	0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x68, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x6d, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4c, 0x6f, 0x67, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x73, 0x70, 0x12, 0x4b, 0x0a, 0x04, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x6c, 0x6f, 0x67, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x6c, 0x6f, 0x67, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x4c, 0x6f, 0x67, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x04, 0x6c, 0x6f,
	0x67, 0x73, 0x22, 0xb5, 0x02, 0x0a, 0x15, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x4c, 0x6f, 0x67, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x14, 0x0a, 0x05,
	0x61, 0x70, 0x70, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x61, 0x70, 0x70,
	0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x65, 0x6e, 0x76, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x12, 0x22, 0x0a, 0x0c, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x65,
	0x66, 0x6f, 0x72, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x41, 0x66, 0x74, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x41, 0x66, 0x74, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x2a, 0x9c, 0x01, 0x0a, 0x0c, 0x4c,
	0x6f, 0x67, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x15, 0x0a, 0x11, 0x4f,
	0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e,
	0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x4f, 0x50, 0x45, 0x52,
	0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x10, 0x02, 0x12, 0x14,
	0x0a, 0x10, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x45, 0x4c, 0x45,
	0x54, 0x45, 0x10, 0x03, 0x12, 0x15, 0x0a, 0x11, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x53, 0x48, 0x10, 0x04, 0x12, 0x1c, 0x0a, 0x18, 0x4f,
	0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x5f,
	0x50, 0x55, 0x42, 0x4c, 0x49, 0x53, 0x48, 0x10, 0x05, 0x32, 0xe6, 0x01, 0x0a, 0x0f, 0x52, 0x65,
	0x6d, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4c, 0x6f, 0x67, 0x12, 0xd2, 0x01,
	0x0a, 0x16, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x4c, 0x6f, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x3b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6c, 0x6f, 0x67, 0x2e, 0x72, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6c, 0x6f, 0x67, 0x2e, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x6d, 0x6f, 0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4c, 0x6f, 0x67, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x3b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6c, 0x6f, 0x67, 0x2e, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x63,
	0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6c, 0x6f, 0x67, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4c, 0x6f, 0x67, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x73, 0x70, 0x22, 0x3e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x38, 0x3a, 0x01, 0x2a, 0x22, 0x33, 0x2f,
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x5f, 0x6c, 0x6f, 0x67, 0x2f, 0x67, 0x65, 0x74, 0x5f, 0x72, 0x65, 0x6d, 0x6f,
	0x74, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x6c, 0x6f, 0x67, 0x5f, 0x6c, 0x69,
	0x73, 0x74, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x61, 0x6d, 0x65, 0x2d, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x6c,
	0x6f, 0x67, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6c,
	0x6f, 0x67, 0x3b, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x6c,
	0x6f, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescOnce sync.Once
	file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescData = file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDesc
)

func file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescGZIP() []byte {
	file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescOnce.Do(func() {
		file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescData)
	})
	return file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDescData
}

var file_api_operationlog_remoteconfiglog_remote_config_log_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_api_operationlog_remoteconfiglog_remote_config_log_proto_goTypes = []interface{}{
	(LogOperation)(0),                 // 0: api.operationlog.remoteconfiglog.LogOperation
	(*GetRemoteConfigLogListReq)(nil), // 1: api.operationlog.remoteconfiglog.GetRemoteConfigLogListReq
	(*GetRemoteConfigLogListRsp)(nil), // 2: api.operationlog.remoteconfiglog.GetRemoteConfigLogListRsp
	(*RemoteConfigLogDetail)(nil),     // 3: api.operationlog.remoteconfiglog.RemoteConfigLogDetail
}
var file_api_operationlog_remoteconfiglog_remote_config_log_proto_depIdxs = []int32{
	3, // 0: api.operationlog.remoteconfiglog.GetRemoteConfigLogListRsp.logs:type_name -> api.operationlog.remoteconfiglog.RemoteConfigLogDetail
	1, // 1: api.operationlog.remoteconfiglog.RemoteConfigLog.GetRemoteConfigLogList:input_type -> api.operationlog.remoteconfiglog.GetRemoteConfigLogListReq
	2, // 2: api.operationlog.remoteconfiglog.RemoteConfigLog.GetRemoteConfigLogList:output_type -> api.operationlog.remoteconfiglog.GetRemoteConfigLogListRsp
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_operationlog_remoteconfiglog_remote_config_log_proto_init() }
func file_api_operationlog_remoteconfiglog_remote_config_log_proto_init() {
	if File_api_operationlog_remoteconfiglog_remote_config_log_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRemoteConfigLogListReq); i {
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
		file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRemoteConfigLogListRsp); i {
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
		file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoteConfigLogDetail); i {
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
			RawDescriptor: file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_operationlog_remoteconfiglog_remote_config_log_proto_goTypes,
		DependencyIndexes: file_api_operationlog_remoteconfiglog_remote_config_log_proto_depIdxs,
		EnumInfos:         file_api_operationlog_remoteconfiglog_remote_config_log_proto_enumTypes,
		MessageInfos:      file_api_operationlog_remoteconfiglog_remote_config_log_proto_msgTypes,
	}.Build()
	File_api_operationlog_remoteconfiglog_remote_config_log_proto = out.File
	file_api_operationlog_remoteconfiglog_remote_config_log_proto_rawDesc = nil
	file_api_operationlog_remoteconfiglog_remote_config_log_proto_goTypes = nil
	file_api_operationlog_remoteconfiglog_remote_config_log_proto_depIdxs = nil
}
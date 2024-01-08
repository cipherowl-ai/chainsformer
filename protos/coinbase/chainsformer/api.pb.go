// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: coinbase/chainsformer/api.proto

package chainsformer

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

type GetFlightInfoCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Query:
	//
	//	*GetFlightInfoCmd_BatchQuery_
	//	*GetFlightInfoCmd_StreamQuery_
	Query isGetFlightInfoCmd_Query `protobuf_oneof:"query"`
}

func (x *GetFlightInfoCmd) Reset() {
	*x = GetFlightInfoCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coinbase_chainsformer_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFlightInfoCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFlightInfoCmd) ProtoMessage() {}

func (x *GetFlightInfoCmd) ProtoReflect() protoreflect.Message {
	mi := &file_coinbase_chainsformer_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFlightInfoCmd.ProtoReflect.Descriptor instead.
func (*GetFlightInfoCmd) Descriptor() ([]byte, []int) {
	return file_coinbase_chainsformer_api_proto_rawDescGZIP(), []int{0}
}

func (m *GetFlightInfoCmd) GetQuery() isGetFlightInfoCmd_Query {
	if m != nil {
		return m.Query
	}
	return nil
}

func (x *GetFlightInfoCmd) GetBatchQuery() *GetFlightInfoCmd_BatchQuery {
	if x, ok := x.GetQuery().(*GetFlightInfoCmd_BatchQuery_); ok {
		return x.BatchQuery
	}
	return nil
}

func (x *GetFlightInfoCmd) GetStreamQuery() *GetFlightInfoCmd_StreamQuery {
	if x, ok := x.GetQuery().(*GetFlightInfoCmd_StreamQuery_); ok {
		return x.StreamQuery
	}
	return nil
}

type isGetFlightInfoCmd_Query interface {
	isGetFlightInfoCmd_Query()
}

type GetFlightInfoCmd_BatchQuery_ struct {
	BatchQuery *GetFlightInfoCmd_BatchQuery `protobuf:"bytes,1,opt,name=batch_query,json=batchQuery,proto3,oneof"`
}

type GetFlightInfoCmd_StreamQuery_ struct {
	StreamQuery *GetFlightInfoCmd_StreamQuery `protobuf:"bytes,2,opt,name=stream_query,json=streamQuery,proto3,oneof"`
}

func (*GetFlightInfoCmd_BatchQuery_) isGetFlightInfoCmd_Query() {}

func (*GetFlightInfoCmd_StreamQuery_) isGetFlightInfoCmd_Query() {}

type GetSchemaCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Table    string `protobuf:"bytes,2,opt,name=table,proto3" json:"table,omitempty"`
	Format   string `protobuf:"bytes,3,opt,name=format,proto3" json:"format,omitempty"`
	Encoding string `protobuf:"bytes,4,opt,name=encoding,proto3" json:"encoding,omitempty"`
}

func (x *GetSchemaCmd) Reset() {
	*x = GetSchemaCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coinbase_chainsformer_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetSchemaCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetSchemaCmd) ProtoMessage() {}

func (x *GetSchemaCmd) ProtoReflect() protoreflect.Message {
	mi := &file_coinbase_chainsformer_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetSchemaCmd.ProtoReflect.Descriptor instead.
func (*GetSchemaCmd) Descriptor() ([]byte, []int) {
	return file_coinbase_chainsformer_api_proto_rawDescGZIP(), []int{1}
}

func (x *GetSchemaCmd) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *GetSchemaCmd) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *GetSchemaCmd) GetEncoding() string {
	if x != nil {
		return x.Encoding
	}
	return ""
}

type GetFlightInfoCmd_BatchQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartHeight        uint64 `protobuf:"varint,2,opt,name=start_height,json=startHeight,proto3" json:"start_height,omitempty"`
	EndHeight          uint64 `protobuf:"varint,3,opt,name=end_height,json=endHeight,proto3" json:"end_height,omitempty"`
	BlocksPerPartition uint64 `protobuf:"varint,4,opt,name=blocks_per_partition,json=blocksPerPartition,proto3" json:"blocks_per_partition,omitempty"`
	BlocksPerRecord    uint64 `protobuf:"varint,5,opt,name=blocks_per_record,json=blocksPerRecord,proto3" json:"blocks_per_record,omitempty"`
	Compression        string `protobuf:"bytes,6,opt,name=compression,proto3" json:"compression,omitempty"`
	Table              string `protobuf:"bytes,7,opt,name=table,proto3" json:"table,omitempty"`
	Format             string `protobuf:"bytes,8,opt,name=format,proto3" json:"format,omitempty"`
	Encoding           string `protobuf:"bytes,9,opt,name=encoding,proto3" json:"encoding,omitempty"`
	PartitionBySize    uint64 `protobuf:"varint,10,opt,name=partition_by_size,json=partitionBySize,proto3" json:"partition_by_size,omitempty"`
}

func (x *GetFlightInfoCmd_BatchQuery) Reset() {
	*x = GetFlightInfoCmd_BatchQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coinbase_chainsformer_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFlightInfoCmd_BatchQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFlightInfoCmd_BatchQuery) ProtoMessage() {}

func (x *GetFlightInfoCmd_BatchQuery) ProtoReflect() protoreflect.Message {
	mi := &file_coinbase_chainsformer_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFlightInfoCmd_BatchQuery.ProtoReflect.Descriptor instead.
func (*GetFlightInfoCmd_BatchQuery) Descriptor() ([]byte, []int) {
	return file_coinbase_chainsformer_api_proto_rawDescGZIP(), []int{0, 0}
}

func (x *GetFlightInfoCmd_BatchQuery) GetStartHeight() uint64 {
	if x != nil {
		return x.StartHeight
	}
	return 0
}

func (x *GetFlightInfoCmd_BatchQuery) GetEndHeight() uint64 {
	if x != nil {
		return x.EndHeight
	}
	return 0
}

func (x *GetFlightInfoCmd_BatchQuery) GetBlocksPerPartition() uint64 {
	if x != nil {
		return x.BlocksPerPartition
	}
	return 0
}

func (x *GetFlightInfoCmd_BatchQuery) GetBlocksPerRecord() uint64 {
	if x != nil {
		return x.BlocksPerRecord
	}
	return 0
}

func (x *GetFlightInfoCmd_BatchQuery) GetCompression() string {
	if x != nil {
		return x.Compression
	}
	return ""
}

func (x *GetFlightInfoCmd_BatchQuery) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *GetFlightInfoCmd_BatchQuery) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *GetFlightInfoCmd_BatchQuery) GetEncoding() string {
	if x != nil {
		return x.Encoding
	}
	return ""
}

func (x *GetFlightInfoCmd_BatchQuery) GetPartitionBySize() uint64 {
	if x != nil {
		return x.PartitionBySize
	}
	return 0
}

type GetFlightInfoCmd_StreamQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartSequence      int64  `protobuf:"varint,2,opt,name=start_sequence,json=startSequence,proto3" json:"start_sequence,omitempty"`
	EndSequence        int64  `protobuf:"varint,3,opt,name=end_sequence,json=endSequence,proto3" json:"end_sequence,omitempty"`
	EventsPerPartition uint64 `protobuf:"varint,4,opt,name=events_per_partition,json=eventsPerPartition,proto3" json:"events_per_partition,omitempty"`
	EventsPerRecord    uint64 `protobuf:"varint,5,opt,name=events_per_record,json=eventsPerRecord,proto3" json:"events_per_record,omitempty"`
	Compression        string `protobuf:"bytes,6,opt,name=compression,proto3" json:"compression,omitempty"`
	Table              string `protobuf:"bytes,7,opt,name=table,proto3" json:"table,omitempty"`
	Format             string `protobuf:"bytes,8,opt,name=format,proto3" json:"format,omitempty"`
	Encoding           string `protobuf:"bytes,9,opt,name=encoding,proto3" json:"encoding,omitempty"`
	PartitionBySize    uint64 `protobuf:"varint,10,opt,name=partition_by_size,json=partitionBySize,proto3" json:"partition_by_size,omitempty"`
}

func (x *GetFlightInfoCmd_StreamQuery) Reset() {
	*x = GetFlightInfoCmd_StreamQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_coinbase_chainsformer_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFlightInfoCmd_StreamQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFlightInfoCmd_StreamQuery) ProtoMessage() {}

func (x *GetFlightInfoCmd_StreamQuery) ProtoReflect() protoreflect.Message {
	mi := &file_coinbase_chainsformer_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFlightInfoCmd_StreamQuery.ProtoReflect.Descriptor instead.
func (*GetFlightInfoCmd_StreamQuery) Descriptor() ([]byte, []int) {
	return file_coinbase_chainsformer_api_proto_rawDescGZIP(), []int{0, 1}
}

func (x *GetFlightInfoCmd_StreamQuery) GetStartSequence() int64 {
	if x != nil {
		return x.StartSequence
	}
	return 0
}

func (x *GetFlightInfoCmd_StreamQuery) GetEndSequence() int64 {
	if x != nil {
		return x.EndSequence
	}
	return 0
}

func (x *GetFlightInfoCmd_StreamQuery) GetEventsPerPartition() uint64 {
	if x != nil {
		return x.EventsPerPartition
	}
	return 0
}

func (x *GetFlightInfoCmd_StreamQuery) GetEventsPerRecord() uint64 {
	if x != nil {
		return x.EventsPerRecord
	}
	return 0
}

func (x *GetFlightInfoCmd_StreamQuery) GetCompression() string {
	if x != nil {
		return x.Compression
	}
	return ""
}

func (x *GetFlightInfoCmd_StreamQuery) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *GetFlightInfoCmd_StreamQuery) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *GetFlightInfoCmd_StreamQuery) GetEncoding() string {
	if x != nil {
		return x.Encoding
	}
	return ""
}

func (x *GetFlightInfoCmd_StreamQuery) GetPartitionBySize() uint64 {
	if x != nil {
		return x.PartitionBySize
	}
	return 0
}

var File_coinbase_chainsformer_api_proto protoreflect.FileDescriptor

var file_coinbase_chainsformer_api_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x63, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x15, 0x63, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x22, 0xef, 0x06, 0x0a, 0x10, 0x47, 0x65, 0x74,
	0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x43, 0x6d, 0x64, 0x12, 0x55, 0x0a,
	0x0b, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x32, 0x2e, 0x63, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x63, 0x68,
	0x61, 0x69, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6c,
	0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x43, 0x6d, 0x64, 0x2e, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x51, 0x75, 0x65, 0x72, 0x79, 0x48, 0x00, 0x52, 0x0a, 0x62, 0x61, 0x74, 0x63, 0x68, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x12, 0x58, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x5f, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x63, 0x6f, 0x69,
	0x6e, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d,
	0x65, 0x72, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x43, 0x6d, 0x64, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x51, 0x75, 0x65, 0x72, 0x79, 0x48,
	0x00, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x51, 0x75, 0x65, 0x72, 0x79, 0x1a, 0xca,
	0x02, 0x0a, 0x0a, 0x42, 0x61, 0x74, 0x63, 0x68, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x21, 0x0a,
	0x0c, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x73, 0x74, 0x61, 0x72, 0x74, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x6e, 0x64, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x65, 0x6e, 0x64, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12,
	0x30, 0x0a, 0x14, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x70, 0x61,
	0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x12, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x50, 0x65, 0x72, 0x50, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x2a, 0x0a, 0x11, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f,
	0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x73, 0x50, 0x65, 0x72, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x20, 0x0a,
	0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1a, 0x0a,
	0x08, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x2a, 0x0a, 0x11, 0x70, 0x61, 0x72,
	0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x62, 0x79, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x0a,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x79, 0x53, 0x69, 0x7a, 0x65, 0x4a, 0x04, 0x08, 0x01, 0x10, 0x02, 0x1a, 0xd3, 0x02, 0x0a, 0x0b,
	0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x25, 0x0a, 0x0e, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x53, 0x65, 0x71, 0x75, 0x65, 0x6e,
	0x63, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x65, 0x6e, 0x64, 0x5f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e,
	0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x65, 0x6e, 0x64, 0x53, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x5f,
	0x70, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x12, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x50, 0x65, 0x72, 0x50, 0x61,
	0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x11, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x73, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x50, 0x65, 0x72, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72,
	0x6d, 0x61, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x12,
	0x2a, 0x0a, 0x11, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x62, 0x79, 0x5f,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x70, 0x61, 0x72, 0x74,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x79, 0x53, 0x69, 0x7a, 0x65, 0x4a, 0x04, 0x08, 0x01, 0x10,
	0x02, 0x42, 0x07, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x22, 0x58, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x43, 0x6d, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x63, 0x6f,
	0x64, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x63, 0x6f,
	0x64, 0x69, 0x6e, 0x67, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x63, 0x68, 0x61, 0x69,
	0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f,
	0x63, 0x6f, 0x69, 0x6e, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x66,
	0x6f, 0x72, 0x6d, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_coinbase_chainsformer_api_proto_rawDescOnce sync.Once
	file_coinbase_chainsformer_api_proto_rawDescData = file_coinbase_chainsformer_api_proto_rawDesc
)

func file_coinbase_chainsformer_api_proto_rawDescGZIP() []byte {
	file_coinbase_chainsformer_api_proto_rawDescOnce.Do(func() {
		file_coinbase_chainsformer_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_coinbase_chainsformer_api_proto_rawDescData)
	})
	return file_coinbase_chainsformer_api_proto_rawDescData
}

var file_coinbase_chainsformer_api_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_coinbase_chainsformer_api_proto_goTypes = []interface{}{
	(*GetFlightInfoCmd)(nil),             // 0: coinbase.chainsformer.GetFlightInfoCmd
	(*GetSchemaCmd)(nil),                 // 1: coinbase.chainsformer.GetSchemaCmd
	(*GetFlightInfoCmd_BatchQuery)(nil),  // 2: coinbase.chainsformer.GetFlightInfoCmd.BatchQuery
	(*GetFlightInfoCmd_StreamQuery)(nil), // 3: coinbase.chainsformer.GetFlightInfoCmd.StreamQuery
}
var file_coinbase_chainsformer_api_proto_depIdxs = []int32{
	2, // 0: coinbase.chainsformer.GetFlightInfoCmd.batch_query:type_name -> coinbase.chainsformer.GetFlightInfoCmd.BatchQuery
	3, // 1: coinbase.chainsformer.GetFlightInfoCmd.stream_query:type_name -> coinbase.chainsformer.GetFlightInfoCmd.StreamQuery
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_coinbase_chainsformer_api_proto_init() }
func file_coinbase_chainsformer_api_proto_init() {
	if File_coinbase_chainsformer_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_coinbase_chainsformer_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFlightInfoCmd); i {
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
		file_coinbase_chainsformer_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetSchemaCmd); i {
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
		file_coinbase_chainsformer_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFlightInfoCmd_BatchQuery); i {
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
		file_coinbase_chainsformer_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFlightInfoCmd_StreamQuery); i {
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
	file_coinbase_chainsformer_api_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*GetFlightInfoCmd_BatchQuery_)(nil),
		(*GetFlightInfoCmd_StreamQuery_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_coinbase_chainsformer_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_coinbase_chainsformer_api_proto_goTypes,
		DependencyIndexes: file_coinbase_chainsformer_api_proto_depIdxs,
		MessageInfos:      file_coinbase_chainsformer_api_proto_msgTypes,
	}.Build()
	File_coinbase_chainsformer_api_proto = out.File
	file_coinbase_chainsformer_api_proto_rawDesc = nil
	file_coinbase_chainsformer_api_proto_goTypes = nil
	file_coinbase_chainsformer_api_proto_depIdxs = nil
}
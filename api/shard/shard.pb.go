// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: shard.proto

package shard

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

// Greet()
type Hello struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShardIndex    uint32                 `protobuf:"varint,1,opt,name=ShardIndex,proto3" json:"ShardIndex,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Hello) Reset() {
	*x = Hello{}
	mi := &file_shard_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Hello) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Hello) ProtoMessage() {}

func (x *Hello) ProtoReflect() protoreflect.Message {
	mi := &file_shard_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Hello.ProtoReflect.Descriptor instead.
func (*Hello) Descriptor() ([]byte, []int) {
	return file_shard_proto_rawDescGZIP(), []int{0}
}

func (x *Hello) GetShardIndex() uint32 {
	if x != nil {
		return x.ShardIndex
	}
	return 0
}

type Hey struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	NumShards      uint32                 `protobuf:"varint,1,opt,name=NumShards,proto3" json:"NumShards,omitempty"`
	ChecksumBase   uint32                 `protobuf:"varint,2,opt,name=ChecksumBase,proto3" json:"ChecksumBase,omitempty"`
	ListenInterval string                 `protobuf:"bytes,3,opt,name=ListenInterval,proto3" json:"ListenInterval,omitempty"`
	SyncInterval   string                 `protobuf:"bytes,4,opt,name=SyncInterval,proto3" json:"SyncInterval,omitempty"`
	// token bucket
	TokenInterval string `protobuf:"bytes,5,opt,name=TokenInterval,proto3" json:"TokenInterval,omitempty"`
	TokensPerHit  uint32 `protobuf:"varint,6,opt,name=TokensPerHit,proto3" json:"TokensPerHit,omitempty"`
	TokensThresh  uint32 `protobuf:"varint,7,opt,name=TokensThresh,proto3" json:"TokensThresh,omitempty"`
	TokensCap     uint32 `protobuf:"varint,8,opt,name=TokensCap,proto3" json:"TokensCap,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Hey) Reset() {
	*x = Hey{}
	mi := &file_shard_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Hey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Hey) ProtoMessage() {}

func (x *Hey) ProtoReflect() protoreflect.Message {
	mi := &file_shard_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Hey.ProtoReflect.Descriptor instead.
func (*Hey) Descriptor() ([]byte, []int) {
	return file_shard_proto_rawDescGZIP(), []int{1}
}

func (x *Hey) GetNumShards() uint32 {
	if x != nil {
		return x.NumShards
	}
	return 0
}

func (x *Hey) GetChecksumBase() uint32 {
	if x != nil {
		return x.ChecksumBase
	}
	return 0
}

func (x *Hey) GetListenInterval() string {
	if x != nil {
		return x.ListenInterval
	}
	return ""
}

func (x *Hey) GetSyncInterval() string {
	if x != nil {
		return x.SyncInterval
	}
	return ""
}

func (x *Hey) GetTokenInterval() string {
	if x != nil {
		return x.TokenInterval
	}
	return ""
}

func (x *Hey) GetTokensPerHit() uint32 {
	if x != nil {
		return x.TokensPerHit
	}
	return 0
}

func (x *Hey) GetTokensThresh() uint32 {
	if x != nil {
		return x.TokensThresh
	}
	return 0
}

func (x *Hey) GetTokensCap() uint32 {
	if x != nil {
		return x.TokensCap
	}
	return 0
}

// SyncDBs()
type SyncSend struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShardIndex    uint32                 `protobuf:"varint,1,opt,name=ShardIndex,proto3" json:"ShardIndex,omitempty"`
	TokenBuckets  map[string]int32       `protobuf:"bytes,2,rep,name=TokenBuckets,proto3" json:"TokenBuckets,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	ServerUIDs    []uint64               `protobuf:"varint,3,rep,packed,name=ServerUIDs,proto3" json:"ServerUIDs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncSend) Reset() {
	*x = SyncSend{}
	mi := &file_shard_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncSend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncSend) ProtoMessage() {}

func (x *SyncSend) ProtoReflect() protoreflect.Message {
	mi := &file_shard_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncSend.ProtoReflect.Descriptor instead.
func (*SyncSend) Descriptor() ([]byte, []int) {
	return file_shard_proto_rawDescGZIP(), []int{2}
}

func (x *SyncSend) GetShardIndex() uint32 {
	if x != nil {
		return x.ShardIndex
	}
	return 0
}

func (x *SyncSend) GetTokenBuckets() map[string]int32 {
	if x != nil {
		return x.TokenBuckets
	}
	return nil
}

func (x *SyncSend) GetServerUIDs() []uint64 {
	if x != nil {
		return x.ServerUIDs
	}
	return nil
}

type SyncReply struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	NeedsBatch    bool                   `protobuf:"varint,1,opt,name=NeedsBatch,proto3" json:"NeedsBatch,omitempty"`
	TokenBuckets  map[string]int32       `protobuf:"bytes,2,rep,name=TokenBuckets,proto3" json:"TokenBuckets,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	ServerUIDs    []uint64               `protobuf:"varint,3,rep,packed,name=ServerUIDs,proto3" json:"ServerUIDs,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SyncReply) Reset() {
	*x = SyncReply{}
	mi := &file_shard_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SyncReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncReply) ProtoMessage() {}

func (x *SyncReply) ProtoReflect() protoreflect.Message {
	mi := &file_shard_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncReply.ProtoReflect.Descriptor instead.
func (*SyncReply) Descriptor() ([]byte, []int) {
	return file_shard_proto_rawDescGZIP(), []int{3}
}

func (x *SyncReply) GetNeedsBatch() bool {
	if x != nil {
		return x.NeedsBatch
	}
	return false
}

func (x *SyncReply) GetTokenBuckets() map[string]int32 {
	if x != nil {
		return x.TokenBuckets
	}
	return nil
}

func (x *SyncReply) GetServerUIDs() []uint64 {
	if x != nil {
		return x.ServerUIDs
	}
	return nil
}

// BatchOut()
type BatchData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ShardIndex    uint32                 `protobuf:"varint,1,opt,name=ShardIndex,proto3" json:"ShardIndex,omitempty"`
	IsLast        bool                   `protobuf:"varint,2,opt,name=IsLast,proto3" json:"IsLast,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BatchData) Reset() {
	*x = BatchData{}
	mi := &file_shard_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BatchData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchData) ProtoMessage() {}

func (x *BatchData) ProtoReflect() protoreflect.Message {
	mi := &file_shard_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchData.ProtoReflect.Descriptor instead.
func (*BatchData) Descriptor() ([]byte, []int) {
	return file_shard_proto_rawDescGZIP(), []int{4}
}

func (x *BatchData) GetShardIndex() uint32 {
	if x != nil {
		return x.ShardIndex
	}
	return 0
}

func (x *BatchData) GetIsLast() bool {
	if x != nil {
		return x.IsLast
	}
	return false
}

var File_shard_proto protoreflect.FileDescriptor

const file_shard_proto_rawDesc = "" +
	"\n" +
	"\vshard.proto\x12\x05shard\x1a\x1bgoogle/protobuf/empty.proto\"'\n" +
	"\x05Hello\x12\x1e\n" +
	"\n" +
	"ShardIndex\x18\x01 \x01(\rR\n" +
	"ShardIndex\"\x9f\x02\n" +
	"\x03Hey\x12\x1c\n" +
	"\tNumShards\x18\x01 \x01(\rR\tNumShards\x12\"\n" +
	"\fChecksumBase\x18\x02 \x01(\rR\fChecksumBase\x12&\n" +
	"\x0eListenInterval\x18\x03 \x01(\tR\x0eListenInterval\x12\"\n" +
	"\fSyncInterval\x18\x04 \x01(\tR\fSyncInterval\x12$\n" +
	"\rTokenInterval\x18\x05 \x01(\tR\rTokenInterval\x12\"\n" +
	"\fTokensPerHit\x18\x06 \x01(\rR\fTokensPerHit\x12\"\n" +
	"\fTokensThresh\x18\a \x01(\rR\fTokensThresh\x12\x1c\n" +
	"\tTokensCap\x18\b \x01(\rR\tTokensCap\"\xd2\x01\n" +
	"\bSyncSend\x12\x1e\n" +
	"\n" +
	"ShardIndex\x18\x01 \x01(\rR\n" +
	"ShardIndex\x12E\n" +
	"\fTokenBuckets\x18\x02 \x03(\v2!.shard.SyncSend.TokenBucketsEntryR\fTokenBuckets\x12\x1e\n" +
	"\n" +
	"ServerUIDs\x18\x03 \x03(\x04R\n" +
	"ServerUIDs\x1a?\n" +
	"\x11TokenBucketsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\x05R\x05value:\x028\x01\"\xd4\x01\n" +
	"\tSyncReply\x12\x1e\n" +
	"\n" +
	"NeedsBatch\x18\x01 \x01(\bR\n" +
	"NeedsBatch\x12F\n" +
	"\fTokenBuckets\x18\x02 \x03(\v2\".shard.SyncReply.TokenBucketsEntryR\fTokenBuckets\x12\x1e\n" +
	"\n" +
	"ServerUIDs\x18\x03 \x03(\x04R\n" +
	"ServerUIDs\x1a?\n" +
	"\x11TokenBucketsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\x05R\x05value:\x028\x01\"C\n" +
	"\tBatchData\x12\x1e\n" +
	"\n" +
	"ShardIndex\x18\x01 \x01(\rR\n" +
	"ShardIndex\x12\x16\n" +
	"\x06IsLast\x18\x02 \x01(\bR\x06IsLast2\x98\x01\n" +
	"\x0fServiceShardAPI\x12!\n" +
	"\x05Greet\x12\f.shard.Hello\x1a\n" +
	".shard.Hey\x12,\n" +
	"\aSyncDBs\x12\x0f.shard.SyncSend\x1a\x10.shard.SyncReply\x124\n" +
	"\bBatchOut\x12\x10.shard.BatchData\x1a\x16.google.protobuf.EmptyB(Z&github.com/PxnPub/pxnMetrics/api/shardb\x06proto3"

var (
	file_shard_proto_rawDescOnce sync.Once
	file_shard_proto_rawDescData []byte
)

func file_shard_proto_rawDescGZIP() []byte {
	file_shard_proto_rawDescOnce.Do(func() {
		file_shard_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_shard_proto_rawDesc), len(file_shard_proto_rawDesc)))
	})
	return file_shard_proto_rawDescData
}

var file_shard_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_shard_proto_goTypes = []any{
	(*Hello)(nil),         // 0: shard.Hello
	(*Hey)(nil),           // 1: shard.Hey
	(*SyncSend)(nil),      // 2: shard.SyncSend
	(*SyncReply)(nil),     // 3: shard.SyncReply
	(*BatchData)(nil),     // 4: shard.BatchData
	nil,                   // 5: shard.SyncSend.TokenBucketsEntry
	nil,                   // 6: shard.SyncReply.TokenBucketsEntry
	(*emptypb.Empty)(nil), // 7: google.protobuf.Empty
}
var file_shard_proto_depIdxs = []int32{
	5, // 0: shard.SyncSend.TokenBuckets:type_name -> shard.SyncSend.TokenBucketsEntry
	6, // 1: shard.SyncReply.TokenBuckets:type_name -> shard.SyncReply.TokenBucketsEntry
	0, // 2: shard.ServiceShardAPI.Greet:input_type -> shard.Hello
	2, // 3: shard.ServiceShardAPI.SyncDBs:input_type -> shard.SyncSend
	4, // 4: shard.ServiceShardAPI.BatchOut:input_type -> shard.BatchData
	1, // 5: shard.ServiceShardAPI.Greet:output_type -> shard.Hey
	3, // 6: shard.ServiceShardAPI.SyncDBs:output_type -> shard.SyncReply
	7, // 7: shard.ServiceShardAPI.BatchOut:output_type -> google.protobuf.Empty
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_shard_proto_init() }
func file_shard_proto_init() {
	if File_shard_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_shard_proto_rawDesc), len(file_shard_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_shard_proto_goTypes,
		DependencyIndexes: file_shard_proto_depIdxs,
		MessageInfos:      file_shard_proto_msgTypes,
	}.Build()
	File_shard_proto = out.File
	file_shard_proto_goTypes = nil
	file_shard_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: pois-service-api.proto

package pb

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

type StatusCode int32

const (
	StatusCode_Success           StatusCode = 0
	StatusCode_Processing        StatusCode = 1
	StatusCode_InvalidParameters StatusCode = 10001
	StatusCode_InvalidPath       StatusCode = 10002
	StatusCode_InternalError     StatusCode = 10003
	StatusCode_OutOfMemory       StatusCode = 10004
	StatusCode_AlgorithmError    StatusCode = 10005
	StatusCode_UnknownError      StatusCode = 10006
)

// Enum value maps for StatusCode.
var (
	StatusCode_name = map[int32]string{
		0:     "Success",
		1:     "Processing",
		10001: "InvalidParameters",
		10002: "InvalidPath",
		10003: "InternalError",
		10004: "OutOfMemory",
		10005: "AlgorithmError",
		10006: "UnknownError",
	}
	StatusCode_value = map[string]int32{
		"Success":           0,
		"Processing":        1,
		"InvalidParameters": 10001,
		"InvalidPath":       10002,
		"InternalError":     10003,
		"OutOfMemory":       10004,
		"AlgorithmError":    10005,
		"UnknownError":      10006,
	}
)

func (x StatusCode) Enum() *StatusCode {
	p := new(StatusCode)
	*p = x
	return p
}

func (x StatusCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StatusCode) Descriptor() protoreflect.EnumDescriptor {
	return file_pois_service_api_proto_enumTypes[0].Descriptor()
}

func (StatusCode) Type() protoreflect.EnumType {
	return &file_pois_service_api_proto_enumTypes[0]
}

func (x StatusCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StatusCode.Descriptor instead.
func (StatusCode) EnumDescriptor() ([]byte, []int) {
	return file_pois_service_api_proto_rawDescGZIP(), []int{0}
}

type Tag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	T       *Tag_T `protobuf:"bytes,1,opt,name=t,proto3" json:"t,omitempty"`
	PhiHash string `protobuf:"bytes,2,opt,name=phi_hash,json=phiHash,proto3" json:"phi_hash,omitempty"`
	Attest  string `protobuf:"bytes,3,opt,name=attest,proto3" json:"attest,omitempty"`
}

func (x *Tag) Reset() {
	*x = Tag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pois_service_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_pois_service_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_pois_service_api_proto_rawDescGZIP(), []int{0}
}

func (x *Tag) GetT() *Tag_T {
	if x != nil {
		return x.T
	}
	return nil
}

func (x *Tag) GetPhiHash() string {
	if x != nil {
		return x.PhiHash
	}
	return ""
}

func (x *Tag) GetAttest() string {
	if x != nil {
		return x.Attest
	}
	return ""
}

type RequestGenTag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileData   []byte `protobuf:"bytes,1,opt,name=file_data,json=fileData,proto3" json:"file_data,omitempty"`
	BlockNum   uint64 `protobuf:"varint,2,opt,name=block_num,json=blockNum,proto3" json:"block_num,omitempty"`
	Name       string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	CustomData string `protobuf:"bytes,4,opt,name=custom_data,json=customData,proto3" json:"custom_data,omitempty"`
}

func (x *RequestGenTag) Reset() {
	*x = RequestGenTag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pois_service_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestGenTag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestGenTag) ProtoMessage() {}

func (x *RequestGenTag) ProtoReflect() protoreflect.Message {
	mi := &file_pois_service_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestGenTag.ProtoReflect.Descriptor instead.
func (*RequestGenTag) Descriptor() ([]byte, []int) {
	return file_pois_service_api_proto_rawDescGZIP(), []int{1}
}

func (x *RequestGenTag) GetFileData() []byte {
	if x != nil {
		return x.FileData
	}
	return nil
}

func (x *RequestGenTag) GetBlockNum() uint64 {
	if x != nil {
		return x.BlockNum
	}
	return 0
}

func (x *RequestGenTag) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RequestGenTag) GetCustomData() string {
	if x != nil {
		return x.CustomData
	}
	return ""
}

type ResponseGenTag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag *Tag `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
}

func (x *ResponseGenTag) Reset() {
	*x = ResponseGenTag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pois_service_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseGenTag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseGenTag) ProtoMessage() {}

func (x *ResponseGenTag) ProtoReflect() protoreflect.Message {
	mi := &file_pois_service_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseGenTag.ProtoReflect.Descriptor instead.
func (*ResponseGenTag) Descriptor() ([]byte, []int) {
	return file_pois_service_api_proto_rawDescGZIP(), []int{2}
}

func (x *ResponseGenTag) GetTag() *Tag {
	if x != nil {
		return x.Tag
	}
	return nil
}

type RequestBatchVerify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AggProof *RequestBatchVerify_BatchVerifyParam `protobuf:"bytes,1,opt,name=agg_proof,json=aggProof,proto3" json:"agg_proof,omitempty"`
	// 38 bytes raw multihash
	PeerId []byte `protobuf:"bytes,2,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	// 32 bytes public key
	MinerPbk []byte `protobuf:"bytes,3,opt,name=miner_pbk,json=minerPbk,proto3" json:"miner_pbk,omitempty"`
	// 64 bytes sign content
	MinerPeerIdSign []byte                     `protobuf:"bytes,4,opt,name=miner_peer_id_sign,json=minerPeerIdSign,proto3" json:"miner_peer_id_sign,omitempty"`
	Qslices         *RequestBatchVerify_Qslice `protobuf:"bytes,5,opt,name=qslices,proto3" json:"qslices,omitempty"`
}

func (x *RequestBatchVerify) Reset() {
	*x = RequestBatchVerify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pois_service_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestBatchVerify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestBatchVerify) ProtoMessage() {}

func (x *RequestBatchVerify) ProtoReflect() protoreflect.Message {
	mi := &file_pois_service_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestBatchVerify.ProtoReflect.Descriptor instead.
func (*RequestBatchVerify) Descriptor() ([]byte, []int) {
	return file_pois_service_api_proto_rawDescGZIP(), []int{3}
}

func (x *RequestBatchVerify) GetAggProof() *RequestBatchVerify_BatchVerifyParam {
	if x != nil {
		return x.AggProof
	}
	return nil
}

func (x *RequestBatchVerify) GetPeerId() []byte {
	if x != nil {
		return x.PeerId
	}
	return nil
}

func (x *RequestBatchVerify) GetMinerPbk() []byte {
	if x != nil {
		return x.MinerPbk
	}
	return nil
}

func (x *RequestBatchVerify) GetMinerPeerIdSign() []byte {
	if x != nil {
		return x.MinerPeerIdSign
	}
	return nil
}

func (x *RequestBatchVerify) GetQslices() *RequestBatchVerify_Qslice {
	if x != nil {
		return x.Qslices
	}
	return nil
}

type ResponseBatchVerify struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BatchVerifyResult  bool     `protobuf:"varint,1,opt,name=batch_verify_result,json=batchVerifyResult,proto3" json:"batch_verify_result,omitempty"`
	TeePeerId          []byte   `protobuf:"bytes,2,opt,name=tee_peer_id,json=teePeerId,proto3" json:"tee_peer_id,omitempty"`
	ServiceBloomFilter []uint64 `protobuf:"varint,3,rep,packed,name=service_bloom_filter,json=serviceBloomFilter,proto3" json:"service_bloom_filter,omitempty"`
	Signature          []byte   `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *ResponseBatchVerify) Reset() {
	*x = ResponseBatchVerify{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pois_service_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseBatchVerify) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseBatchVerify) ProtoMessage() {}

func (x *ResponseBatchVerify) ProtoReflect() protoreflect.Message {
	mi := &file_pois_service_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseBatchVerify.ProtoReflect.Descriptor instead.
func (*ResponseBatchVerify) Descriptor() ([]byte, []int) {
	return file_pois_service_api_proto_rawDescGZIP(), []int{4}
}

func (x *ResponseBatchVerify) GetBatchVerifyResult() bool {
	if x != nil {
		return x.BatchVerifyResult
	}
	return false
}

func (x *ResponseBatchVerify) GetTeePeerId() []byte {
	if x != nil {
		return x.TeePeerId
	}
	return nil
}

func (x *ResponseBatchVerify) GetServiceBloomFilter() []uint64 {
	if x != nil {
		return x.ServiceBloomFilter
	}
	return nil
}

func (x *ResponseBatchVerify) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type Tag_T struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	U    string   `protobuf:"bytes,2,opt,name=u,proto3" json:"u,omitempty"`
	Phi  []string `protobuf:"bytes,3,rep,name=phi,proto3" json:"phi,omitempty"`
}

func (x *Tag_T) Reset() {
	*x = Tag_T{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pois_service_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag_T) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag_T) ProtoMessage() {}

func (x *Tag_T) ProtoReflect() protoreflect.Message {
	mi := &file_pois_service_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag_T.ProtoReflect.Descriptor instead.
func (*Tag_T) Descriptor() ([]byte, []int) {
	return file_pois_service_api_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Tag_T) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tag_T) GetU() string {
	if x != nil {
		return x.U
	}
	return ""
}

func (x *Tag_T) GetPhi() []string {
	if x != nil {
		return x.Phi
	}
	return nil
}

type RequestBatchVerify_Qslice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RandomIndexList []uint32 `protobuf:"varint,1,rep,packed,name=random_index_list,json=randomIndexList,proto3" json:"random_index_list,omitempty"`
	RandomList      [][]byte `protobuf:"bytes,2,rep,name=random_list,json=randomList,proto3" json:"random_list,omitempty"`
}

func (x *RequestBatchVerify_Qslice) Reset() {
	*x = RequestBatchVerify_Qslice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pois_service_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestBatchVerify_Qslice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestBatchVerify_Qslice) ProtoMessage() {}

func (x *RequestBatchVerify_Qslice) ProtoReflect() protoreflect.Message {
	mi := &file_pois_service_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestBatchVerify_Qslice.ProtoReflect.Descriptor instead.
func (*RequestBatchVerify_Qslice) Descriptor() ([]byte, []int) {
	return file_pois_service_api_proto_rawDescGZIP(), []int{3, 0}
}

func (x *RequestBatchVerify_Qslice) GetRandomIndexList() []uint32 {
	if x != nil {
		return x.RandomIndexList
	}
	return nil
}

func (x *RequestBatchVerify_Qslice) GetRandomList() [][]byte {
	if x != nil {
		return x.RandomList
	}
	return nil
}

type RequestBatchVerify_BatchVerifyParam struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Names []string `protobuf:"bytes,1,rep,name=names,proto3" json:"names,omitempty"`
	Us    []string `protobuf:"bytes,2,rep,name=us,proto3" json:"us,omitempty"`
	Mus   []string `protobuf:"bytes,3,rep,name=mus,proto3" json:"mus,omitempty"`
	Sigma string   `protobuf:"bytes,4,opt,name=sigma,proto3" json:"sigma,omitempty"`
}

func (x *RequestBatchVerify_BatchVerifyParam) Reset() {
	*x = RequestBatchVerify_BatchVerifyParam{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pois_service_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestBatchVerify_BatchVerifyParam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestBatchVerify_BatchVerifyParam) ProtoMessage() {}

func (x *RequestBatchVerify_BatchVerifyParam) ProtoReflect() protoreflect.Message {
	mi := &file_pois_service_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestBatchVerify_BatchVerifyParam.ProtoReflect.Descriptor instead.
func (*RequestBatchVerify_BatchVerifyParam) Descriptor() ([]byte, []int) {
	return file_pois_service_api_proto_rawDescGZIP(), []int{3, 1}
}

func (x *RequestBatchVerify_BatchVerifyParam) GetNames() []string {
	if x != nil {
		return x.Names
	}
	return nil
}

func (x *RequestBatchVerify_BatchVerifyParam) GetUs() []string {
	if x != nil {
		return x.Us
	}
	return nil
}

func (x *RequestBatchVerify_BatchVerifyParam) GetMus() []string {
	if x != nil {
		return x.Mus
	}
	return nil
}

func (x *RequestBatchVerify_BatchVerifyParam) GetSigma() string {
	if x != nil {
		return x.Sigma
	}
	return ""
}

var File_pois_service_api_proto protoreflect.FileDescriptor

var file_pois_service_api_proto_rawDesc = []byte{
	0x0a, 0x16, 0x70, 0x6f, 0x69, 0x73, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2d, 0x61,
	0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x87, 0x01, 0x0a, 0x03, 0x54, 0x61, 0x67,
	0x12, 0x14, 0x0a, 0x01, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x54, 0x61,
	0x67, 0x2e, 0x54, 0x52, 0x01, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x68, 0x69, 0x5f, 0x68, 0x61,
	0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x68, 0x69, 0x48, 0x61, 0x73,
	0x68, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x1a, 0x37, 0x0a, 0x01, 0x54, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x75, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x01, 0x75,
	0x12, 0x10, 0x0a, 0x03, 0x70, 0x68, 0x69, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x70,
	0x68, 0x69, 0x22, 0x7e, 0x0a, 0x0d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x47, 0x65, 0x6e,
	0x54, 0x61, 0x67, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1b, 0x0a, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x08, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x4e, 0x75, 0x6d, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x44, 0x61,
	0x74, 0x61, 0x22, 0x28, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x47, 0x65,
	0x6e, 0x54, 0x61, 0x67, 0x12, 0x16, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x04, 0x2e, 0x54, 0x61, 0x67, 0x52, 0x03, 0x74, 0x61, 0x67, 0x22, 0xa9, 0x03, 0x0a,
	0x12, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x56, 0x65, 0x72,
	0x69, 0x66, 0x79, 0x12, 0x41, 0x0a, 0x09, 0x61, 0x67, 0x67, 0x5f, 0x70, 0x72, 0x6f, 0x6f, 0x66,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x2e, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x52, 0x08, 0x61, 0x67,
	0x67, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1b, 0x0a, 0x09, 0x6d, 0x69, 0x6e, 0x65, 0x72, 0x5f, 0x70, 0x62, 0x6b, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x08, 0x6d, 0x69, 0x6e, 0x65, 0x72, 0x50, 0x62, 0x6b, 0x12, 0x2b, 0x0a, 0x12,
	0x6d, 0x69, 0x6e, 0x65, 0x72, 0x5f, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x5f, 0x73, 0x69,
	0x67, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x6d, 0x69, 0x6e, 0x65, 0x72, 0x50,
	0x65, 0x65, 0x72, 0x49, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x12, 0x34, 0x0a, 0x07, 0x71, 0x73, 0x6c,
	0x69, 0x63, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x2e,
	0x51, 0x73, 0x6c, 0x69, 0x63, 0x65, 0x52, 0x07, 0x71, 0x73, 0x6c, 0x69, 0x63, 0x65, 0x73, 0x1a,
	0x55, 0x0a, 0x06, 0x51, 0x73, 0x6c, 0x69, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x72, 0x61, 0x6e,
	0x64, 0x6f, 0x6d, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0d, 0x52, 0x0f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x5f,
	0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0a, 0x72, 0x61, 0x6e, 0x64,
	0x6f, 0x6d, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x60, 0x0a, 0x10, 0x42, 0x61, 0x74, 0x63, 0x68, 0x56,
	0x65, 0x72, 0x69, 0x66, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x12, 0x0e, 0x0a, 0x02, 0x75, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x02, 0x75, 0x73,
	0x12, 0x10, 0x0a, 0x03, 0x6d, 0x75, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x6d,
	0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x69, 0x67, 0x6d, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x73, 0x69, 0x67, 0x6d, 0x61, 0x22, 0xb5, 0x01, 0x0a, 0x13, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x12, 0x2e, 0x0a, 0x13, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79,
	0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x11, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x12, 0x1e, 0x0a, 0x0b, 0x74, 0x65, 0x65, 0x5f, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x74, 0x65, 0x65, 0x50, 0x65, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x30, 0x0a, 0x14, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x62, 0x6c, 0x6f, 0x6f,
	0x6d, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x03, 0x28, 0x04, 0x52, 0x12,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x2a, 0xa1, 0x01, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x0b, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x10, 0x01, 0x12, 0x16, 0x0a, 0x11,
	0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72,
	0x73, 0x10, 0x91, 0x4e, 0x12, 0x10, 0x0a, 0x0b, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x50,
	0x61, 0x74, 0x68, 0x10, 0x92, 0x4e, 0x12, 0x12, 0x0a, 0x0d, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x93, 0x4e, 0x12, 0x10, 0x0a, 0x0b, 0x4f, 0x75,
	0x74, 0x4f, 0x66, 0x4d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x10, 0x94, 0x4e, 0x12, 0x13, 0x0a, 0x0e,
	0x41, 0x6c, 0x67, 0x6f, 0x72, 0x69, 0x74, 0x68, 0x6d, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x95,
	0x4e, 0x12, 0x11, 0x0a, 0x0c, 0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x10, 0x96, 0x4e, 0x32, 0x85, 0x01, 0x0a, 0x08, 0x50, 0x6f, 0x64, 0x72, 0x32, 0x41, 0x70,
	0x69, 0x12, 0x34, 0x0a, 0x0f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x67, 0x65, 0x6e,
	0x5f, 0x74, 0x61, 0x67, 0x12, 0x0e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x47, 0x65,
	0x6e, 0x54, 0x61, 0x67, 0x1a, 0x0f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x47,
	0x65, 0x6e, 0x54, 0x61, 0x67, 0x22, 0x00, 0x12, 0x43, 0x0a, 0x14, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x5f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x76, 0x65, 0x72, 0x69, 0x66, 0x79, 0x12,
	0x13, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x1a, 0x14, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x22, 0x00, 0x42, 0x07, 0x5a, 0x05,
	0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pois_service_api_proto_rawDescOnce sync.Once
	file_pois_service_api_proto_rawDescData = file_pois_service_api_proto_rawDesc
)

func file_pois_service_api_proto_rawDescGZIP() []byte {
	file_pois_service_api_proto_rawDescOnce.Do(func() {
		file_pois_service_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_pois_service_api_proto_rawDescData)
	})
	return file_pois_service_api_proto_rawDescData
}

var file_pois_service_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pois_service_api_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_pois_service_api_proto_goTypes = []interface{}{
	(StatusCode)(0),                             // 0: StatusCode
	(*Tag)(nil),                                 // 1: Tag
	(*RequestGenTag)(nil),                       // 2: RequestGenTag
	(*ResponseGenTag)(nil),                      // 3: ResponseGenTag
	(*RequestBatchVerify)(nil),                  // 4: RequestBatchVerify
	(*ResponseBatchVerify)(nil),                 // 5: ResponseBatchVerify
	(*Tag_T)(nil),                               // 6: Tag.T
	(*RequestBatchVerify_Qslice)(nil),           // 7: RequestBatchVerify.Qslice
	(*RequestBatchVerify_BatchVerifyParam)(nil), // 8: RequestBatchVerify.BatchVerifyParam
}
var file_pois_service_api_proto_depIdxs = []int32{
	6, // 0: Tag.t:type_name -> Tag.T
	1, // 1: ResponseGenTag.tag:type_name -> Tag
	8, // 2: RequestBatchVerify.agg_proof:type_name -> RequestBatchVerify.BatchVerifyParam
	7, // 3: RequestBatchVerify.qslices:type_name -> RequestBatchVerify.Qslice
	2, // 4: Podr2Api.request_gen_tag:input_type -> RequestGenTag
	4, // 5: Podr2Api.request_batch_verify:input_type -> RequestBatchVerify
	3, // 6: Podr2Api.request_gen_tag:output_type -> ResponseGenTag
	5, // 7: Podr2Api.request_batch_verify:output_type -> ResponseBatchVerify
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_pois_service_api_proto_init() }
func file_pois_service_api_proto_init() {
	if File_pois_service_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pois_service_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tag); i {
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
		file_pois_service_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestGenTag); i {
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
		file_pois_service_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseGenTag); i {
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
		file_pois_service_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestBatchVerify); i {
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
		file_pois_service_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseBatchVerify); i {
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
		file_pois_service_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tag_T); i {
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
		file_pois_service_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestBatchVerify_Qslice); i {
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
		file_pois_service_api_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestBatchVerify_BatchVerifyParam); i {
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
			RawDescriptor: file_pois_service_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pois_service_api_proto_goTypes,
		DependencyIndexes: file_pois_service_api_proto_depIdxs,
		EnumInfos:         file_pois_service_api_proto_enumTypes,
		MessageInfos:      file_pois_service_api_proto_msgTypes,
	}.Build()
	File_pois_service_api_proto = out.File
	file_pois_service_api_proto_rawDesc = nil
	file_pois_service_api_proto_goTypes = nil
	file_pois_service_api_proto_depIdxs = nil
}

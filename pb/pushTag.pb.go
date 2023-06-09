// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: pushTag.proto

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

type T struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	U    string   `protobuf:"bytes,2,opt,name=u,proto3" json:"u,omitempty"`
	Phi  []string `protobuf:"bytes,3,rep,name=phi,proto3" json:"phi,omitempty"`
}

func (x *T) Reset() {
	*x = T{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pushTag_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *T) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*T) ProtoMessage() {}

func (x *T) ProtoReflect() protoreflect.Message {
	mi := &file_pushTag_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use T.ProtoReflect.Descriptor instead.
func (*T) Descriptor() ([]byte, []int) {
	return file_pushTag_proto_rawDescGZIP(), []int{0}
}

func (x *T) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *T) GetU() string {
	if x != nil {
		return x.U
	}
	return ""
}

func (x *T) GetPhi() []string {
	if x != nil {
		return x.Phi
	}
	return nil
}

type Tag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	T       *T     `protobuf:"bytes,1,opt,name=t,proto3" json:"t,omitempty"`
	PhiHash string `protobuf:"bytes,2,opt,name=phi_hash,json=phiHash,proto3" json:"phi_hash,omitempty"`
	Attest  string `protobuf:"bytes,3,opt,name=attest,proto3" json:"attest,omitempty"`
}

func (x *Tag) Reset() {
	*x = Tag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pushTag_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_pushTag_proto_msgTypes[1]
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
	return file_pushTag_proto_rawDescGZIP(), []int{1}
}

func (x *Tag) GetT() *T {
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

type CustomTagGenResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag *Tag `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
}

func (x *CustomTagGenResult) Reset() {
	*x = CustomTagGenResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pushTag_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CustomTagGenResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CustomTagGenResult) ProtoMessage() {}

func (x *CustomTagGenResult) ProtoReflect() protoreflect.Message {
	mi := &file_pushTag_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CustomTagGenResult.ProtoReflect.Descriptor instead.
func (*CustomTagGenResult) Descriptor() ([]byte, []int) {
	return file_pushTag_proto_rawDescGZIP(), []int{2}
}

func (x *CustomTagGenResult) GetTag() *Tag {
	if x != nil {
		return x.Tag
	}
	return nil
}

type IdleTagGenResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tag  *Tag   `protobuf:"bytes,1,opt,name=tag,proto3" json:"tag,omitempty"`
	Sign []byte `protobuf:"bytes,2,opt,name=sign,proto3" json:"sign,omitempty"`
}

func (x *IdleTagGenResult) Reset() {
	*x = IdleTagGenResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pushTag_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdleTagGenResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdleTagGenResult) ProtoMessage() {}

func (x *IdleTagGenResult) ProtoReflect() protoreflect.Message {
	mi := &file_pushTag_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdleTagGenResult.ProtoReflect.Descriptor instead.
func (*IdleTagGenResult) Descriptor() ([]byte, []int) {
	return file_pushTag_proto_rawDescGZIP(), []int{3}
}

func (x *IdleTagGenResult) GetTag() *Tag {
	if x != nil {
		return x.Tag
	}
	return nil
}

func (x *IdleTagGenResult) GetSign() []byte {
	if x != nil {
		return x.Sign
	}
	return nil
}

type GenErrorResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *GenErrorResult) Reset() {
	*x = GenErrorResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pushTag_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenErrorResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenErrorResult) ProtoMessage() {}

func (x *GenErrorResult) ProtoReflect() protoreflect.Message {
	mi := &file_pushTag_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenErrorResult.ProtoReflect.Descriptor instead.
func (*GenErrorResult) Descriptor() ([]byte, []int) {
	return file_pushTag_proto_rawDescGZIP(), []int{4}
}

func (x *GenErrorResult) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *GenErrorResult) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type TagPushRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Result:
	//	*TagPushRequest_Ctgr
	//	*TagPushRequest_Itgr
	//	*TagPushRequest_Error
	Result isTagPushRequest_Result `protobuf_oneof:"result"`
}

func (x *TagPushRequest) Reset() {
	*x = TagPushRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pushTag_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TagPushRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TagPushRequest) ProtoMessage() {}

func (x *TagPushRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pushTag_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TagPushRequest.ProtoReflect.Descriptor instead.
func (*TagPushRequest) Descriptor() ([]byte, []int) {
	return file_pushTag_proto_rawDescGZIP(), []int{5}
}

func (m *TagPushRequest) GetResult() isTagPushRequest_Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (x *TagPushRequest) GetCtgr() *CustomTagGenResult {
	if x, ok := x.GetResult().(*TagPushRequest_Ctgr); ok {
		return x.Ctgr
	}
	return nil
}

func (x *TagPushRequest) GetItgr() *IdleTagGenResult {
	if x, ok := x.GetResult().(*TagPushRequest_Itgr); ok {
		return x.Itgr
	}
	return nil
}

func (x *TagPushRequest) GetError() *GenErrorResult {
	if x, ok := x.GetResult().(*TagPushRequest_Error); ok {
		return x.Error
	}
	return nil
}

type isTagPushRequest_Result interface {
	isTagPushRequest_Result()
}

type TagPushRequest_Ctgr struct {
	Ctgr *CustomTagGenResult `protobuf:"bytes,1,opt,name=ctgr,proto3,oneof"`
}

type TagPushRequest_Itgr struct {
	Itgr *IdleTagGenResult `protobuf:"bytes,2,opt,name=itgr,proto3,oneof"`
}

type TagPushRequest_Error struct {
	Error *GenErrorResult `protobuf:"bytes,3,opt,name=error,proto3,oneof"`
}

func (*TagPushRequest_Ctgr) isTagPushRequest_Result() {}

func (*TagPushRequest_Itgr) isTagPushRequest_Result() {}

func (*TagPushRequest_Error) isTagPushRequest_Result() {}

type TagPushResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *TagPushResponse) Reset() {
	*x = TagPushResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pushTag_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TagPushResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TagPushResponse) ProtoMessage() {}

func (x *TagPushResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pushTag_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TagPushResponse.ProtoReflect.Descriptor instead.
func (*TagPushResponse) Descriptor() ([]byte, []int) {
	return file_pushTag_proto_rawDescGZIP(), []int{6}
}

func (x *TagPushResponse) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

var File_pushTag_proto protoreflect.FileDescriptor

var file_pushTag_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x75, 0x73, 0x68, 0x54, 0x61, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x37, 0x0a, 0x01, 0x54, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x0c, 0x0a, 0x01, 0x75, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x01, 0x75, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x68, 0x69, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x03, 0x70, 0x68, 0x69, 0x22, 0x4a, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x12,
	0x10, 0x0a, 0x01, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x02, 0x2e, 0x54, 0x52, 0x01,
	0x74, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x68, 0x69, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x68, 0x69, 0x48, 0x61, 0x73, 0x68, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x74,
	0x74, 0x65, 0x73, 0x74, 0x22, 0x2c, 0x0a, 0x12, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x54, 0x61,
	0x67, 0x47, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x03, 0x74, 0x61,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x54, 0x61, 0x67, 0x52, 0x03, 0x74,
	0x61, 0x67, 0x22, 0x3e, 0x0a, 0x10, 0x49, 0x64, 0x6c, 0x65, 0x54, 0x61, 0x67, 0x47, 0x65, 0x6e,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x54, 0x61, 0x67, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x69, 0x67, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x73, 0x69,
	0x67, 0x6e, 0x22, 0x36, 0x0a, 0x0e, 0x47, 0x65, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x97, 0x01, 0x0a, 0x0e, 0x54,
	0x61, 0x67, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x29, 0x0a,
	0x04, 0x63, 0x74, 0x67, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x43, 0x75,
	0x73, 0x74, 0x6f, 0x6d, 0x54, 0x61, 0x67, 0x47, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x48, 0x00, 0x52, 0x04, 0x63, 0x74, 0x67, 0x72, 0x12, 0x27, 0x0a, 0x04, 0x69, 0x74, 0x67, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x49, 0x64, 0x6c, 0x65, 0x54, 0x61, 0x67,
	0x47, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x48, 0x00, 0x52, 0x04, 0x69, 0x74, 0x67,
	0x72, 0x12, 0x27, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0f, 0x2e, 0x47, 0x65, 0x6e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x48, 0x00, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x42, 0x08, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x22, 0x25, 0x0a, 0x0f, 0x54, 0x61, 0x67, 0x50, 0x75, 0x73, 0x68, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x42, 0x07, 0x5a, 0x05, 0x2e,
	0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pushTag_proto_rawDescOnce sync.Once
	file_pushTag_proto_rawDescData = file_pushTag_proto_rawDesc
)

func file_pushTag_proto_rawDescGZIP() []byte {
	file_pushTag_proto_rawDescOnce.Do(func() {
		file_pushTag_proto_rawDescData = protoimpl.X.CompressGZIP(file_pushTag_proto_rawDescData)
	})
	return file_pushTag_proto_rawDescData
}

var file_pushTag_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pushTag_proto_goTypes = []interface{}{
	(*T)(nil),                  // 0: T
	(*Tag)(nil),                // 1: Tag
	(*CustomTagGenResult)(nil), // 2: CustomTagGenResult
	(*IdleTagGenResult)(nil),   // 3: IdleTagGenResult
	(*GenErrorResult)(nil),     // 4: GenErrorResult
	(*TagPushRequest)(nil),     // 5: TagPushRequest
	(*TagPushResponse)(nil),    // 6: TagPushResponse
}
var file_pushTag_proto_depIdxs = []int32{
	0, // 0: Tag.t:type_name -> T
	1, // 1: CustomTagGenResult.tag:type_name -> Tag
	1, // 2: IdleTagGenResult.tag:type_name -> Tag
	2, // 3: TagPushRequest.ctgr:type_name -> CustomTagGenResult
	3, // 4: TagPushRequest.itgr:type_name -> IdleTagGenResult
	4, // 5: TagPushRequest.error:type_name -> GenErrorResult
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_pushTag_proto_init() }
func file_pushTag_proto_init() {
	if File_pushTag_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pushTag_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*T); i {
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
		file_pushTag_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_pushTag_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CustomTagGenResult); i {
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
		file_pushTag_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdleTagGenResult); i {
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
		file_pushTag_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenErrorResult); i {
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
		file_pushTag_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TagPushRequest); i {
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
		file_pushTag_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TagPushResponse); i {
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
	file_pushTag_proto_msgTypes[5].OneofWrappers = []interface{}{
		(*TagPushRequest_Ctgr)(nil),
		(*TagPushRequest_Itgr)(nil),
		(*TagPushRequest_Error)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pushTag_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pushTag_proto_goTypes,
		DependencyIndexes: file_pushTag_proto_depIdxs,
		MessageInfos:      file_pushTag_proto_msgTypes,
	}.Build()
	File_pushTag_proto = out.File
	file_pushTag_proto_rawDesc = nil
	file_pushTag_proto_goTypes = nil
	file_pushTag_proto_depIdxs = nil
}

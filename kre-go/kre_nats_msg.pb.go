// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: kre_nats_msg.proto

package kre

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type KreNatsMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TrackingId string                     `protobuf:"bytes,1,opt,name=tracking_id,json=trackingId,proto3" json:"tracking_id,omitempty"`
	Tracking   []*KreNatsMessage_Tracking `protobuf:"bytes,2,rep,name=tracking,proto3" json:"tracking,omitempty"`
	Reply      string                     `protobuf:"bytes,3,opt,name=reply,proto3" json:"reply,omitempty"`
	Payload    *anypb.Any                 `protobuf:"bytes,4,opt,name=payload,proto3" json:"payload,omitempty"`
	Error      string                     `protobuf:"bytes,5,opt,name=error,proto3" json:"error,omitempty"`
	Replied    bool                       `protobuf:"varint,6,opt,name=replied,proto3" json:"replied,omitempty"`
}

func (x *KreNatsMessage) Reset() {
	*x = KreNatsMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kre_nats_msg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KreNatsMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KreNatsMessage) ProtoMessage() {}

func (x *KreNatsMessage) ProtoReflect() protoreflect.Message {
	mi := &file_kre_nats_msg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KreNatsMessage.ProtoReflect.Descriptor instead.
func (*KreNatsMessage) Descriptor() ([]byte, []int) {
	return file_kre_nats_msg_proto_rawDescGZIP(), []int{0}
}

func (x *KreNatsMessage) GetTrackingId() string {
	if x != nil {
		return x.TrackingId
	}
	return ""
}

func (x *KreNatsMessage) GetTracking() []*KreNatsMessage_Tracking {
	if x != nil {
		return x.Tracking
	}
	return nil
}

func (x *KreNatsMessage) GetReply() string {
	if x != nil {
		return x.Reply
	}
	return ""
}

func (x *KreNatsMessage) GetPayload() *anypb.Any {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *KreNatsMessage) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *KreNatsMessage) GetReplied() bool {
	if x != nil {
		return x.Replied
	}
	return false
}

type KreNatsMessage_Tracking struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeName string `protobuf:"bytes,1,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	Start    string `protobuf:"bytes,2,opt,name=start,proto3" json:"start,omitempty"`
	End      string `protobuf:"bytes,3,opt,name=end,proto3" json:"end,omitempty"`
}

func (x *KreNatsMessage_Tracking) Reset() {
	*x = KreNatsMessage_Tracking{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kre_nats_msg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KreNatsMessage_Tracking) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KreNatsMessage_Tracking) ProtoMessage() {}

func (x *KreNatsMessage_Tracking) ProtoReflect() protoreflect.Message {
	mi := &file_kre_nats_msg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KreNatsMessage_Tracking.ProtoReflect.Descriptor instead.
func (*KreNatsMessage_Tracking) Descriptor() ([]byte, []int) {
	return file_kre_nats_msg_proto_rawDescGZIP(), []int{0, 0}
}

func (x *KreNatsMessage_Tracking) GetNodeName() string {
	if x != nil {
		return x.NodeName
	}
	return ""
}

func (x *KreNatsMessage_Tracking) GetStart() string {
	if x != nil {
		return x.Start
	}
	return ""
}

func (x *KreNatsMessage_Tracking) GetEnd() string {
	if x != nil {
		return x.End
	}
	return ""
}

var File_kre_nats_msg_proto protoreflect.FileDescriptor

var file_kre_nats_msg_proto_rawDesc = []byte{
	0x0a, 0x12, 0x6b, 0x72, 0x65, 0x5f, 0x6e, 0x61, 0x74, 0x73, 0x5f, 0x6d, 0x73, 0x67, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xae, 0x02, 0x0a, 0x0e, 0x4b, 0x72, 0x65, 0x4e, 0x61, 0x74, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e,
	0x67, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x4b, 0x72, 0x65, 0x4e, 0x61, 0x74, 0x73, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x52,
	0x08, 0x74, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x70,
	0x6c, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x2e, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x65, 0x64,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x65, 0x64, 0x1a,
	0x4f, 0x0a, 0x08, 0x54, 0x72, 0x61, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x12, 0x1b, 0x0a, 0x09, 0x6e,
	0x6f, 0x64, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x6e, 0x6f, 0x64, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x6e, 0x64,
	0x42, 0x05, 0x5a, 0x03, 0x6b, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kre_nats_msg_proto_rawDescOnce sync.Once
	file_kre_nats_msg_proto_rawDescData = file_kre_nats_msg_proto_rawDesc
)

func file_kre_nats_msg_proto_rawDescGZIP() []byte {
	file_kre_nats_msg_proto_rawDescOnce.Do(func() {
		file_kre_nats_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_kre_nats_msg_proto_rawDescData)
	})
	return file_kre_nats_msg_proto_rawDescData
}

var file_kre_nats_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_kre_nats_msg_proto_goTypes = []interface{}{
	(*KreNatsMessage)(nil),          // 0: KreNatsMessage
	(*KreNatsMessage_Tracking)(nil), // 1: KreNatsMessage.Tracking
	(*anypb.Any)(nil),               // 2: google.protobuf.Any
}
var file_kre_nats_msg_proto_depIdxs = []int32{
	1, // 0: KreNatsMessage.tracking:type_name -> KreNatsMessage.Tracking
	2, // 1: KreNatsMessage.payload:type_name -> google.protobuf.Any
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_kre_nats_msg_proto_init() }
func file_kre_nats_msg_proto_init() {
	if File_kre_nats_msg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kre_nats_msg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KreNatsMessage); i {
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
		file_kre_nats_msg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KreNatsMessage_Tracking); i {
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
			RawDescriptor: file_kre_nats_msg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_kre_nats_msg_proto_goTypes,
		DependencyIndexes: file_kre_nats_msg_proto_depIdxs,
		MessageInfos:      file_kre_nats_msg_proto_msgTypes,
	}.Build()
	File_kre_nats_msg_proto = out.File
	file_kre_nats_msg_proto_rawDesc = nil
	file_kre_nats_msg_proto_goTypes = nil
	file_kre_nats_msg_proto_depIdxs = nil
}

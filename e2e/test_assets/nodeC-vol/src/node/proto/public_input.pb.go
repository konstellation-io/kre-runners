// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.10
// source: public_input.proto

package proto

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

type Testing struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsEarlyReply bool `protobuf:"varint,1,opt,name=is_early_reply,json=isEarlyReply,proto3" json:"is_early_reply,omitempty"`
	IsEarlyExit  bool `protobuf:"varint,2,opt,name=is_early_exit,json=isEarlyExit,proto3" json:"is_early_exit,omitempty"`
	TestStores   bool `protobuf:"varint,3,opt,name=test_stores,json=testStores,proto3" json:"test_stores,omitempty"`
}

func (x *Testing) Reset() {
	*x = Testing{}
	if protoimpl.UnsafeEnabled {
		mi := &file_public_input_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Testing) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Testing) ProtoMessage() {}

func (x *Testing) ProtoReflect() protoreflect.Message {
	mi := &file_public_input_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Testing.ProtoReflect.Descriptor instead.
func (*Testing) Descriptor() ([]byte, []int) {
	return file_public_input_proto_rawDescGZIP(), []int{0}
}

func (x *Testing) GetIsEarlyReply() bool {
	if x != nil {
		return x.IsEarlyReply
	}
	return false
}

func (x *Testing) GetIsEarlyExit() bool {
	if x != nil {
		return x.IsEarlyExit
	}
	return false
}

func (x *Testing) GetTestStores() bool {
	if x != nil {
		return x.TestStores
	}
	return false
}

type TestingResults struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TestStoresSuccess bool `protobuf:"varint,1,opt,name=test_stores_success,json=testStoresSuccess,proto3" json:"test_stores_success,omitempty"`
}

func (x *TestingResults) Reset() {
	*x = TestingResults{}
	if protoimpl.UnsafeEnabled {
		mi := &file_public_input_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestingResults) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestingResults) ProtoMessage() {}

func (x *TestingResults) ProtoReflect() protoreflect.Message {
	mi := &file_public_input_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestingResults.ProtoReflect.Descriptor instead.
func (*TestingResults) Descriptor() ([]byte, []int) {
	return file_public_input_proto_rawDescGZIP(), []int{1}
}

func (x *TestingResults) GetTestStoresSuccess() bool {
	if x != nil {
		return x.TestStoresSuccess
	}
	return false
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Testing *Testing `protobuf:"bytes,2,opt,name=testing,proto3" json:"testing,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_public_input_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_public_input_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_public_input_proto_rawDescGZIP(), []int{2}
}

func (x *Request) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Request) GetTesting() *Testing {
	if x != nil {
		return x.Testing
	}
	return nil
}

type NodeBRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Greeting       string          `protobuf:"bytes,1,opt,name=greeting,proto3" json:"greeting,omitempty"`
	Testing        *Testing        `protobuf:"bytes,2,opt,name=testing,proto3" json:"testing,omitempty"`
	TestingResults *TestingResults `protobuf:"bytes,3,opt,name=testing_results,json=testingResults,proto3" json:"testing_results,omitempty"`
}

func (x *NodeBRequest) Reset() {
	*x = NodeBRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_public_input_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeBRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeBRequest) ProtoMessage() {}

func (x *NodeBRequest) ProtoReflect() protoreflect.Message {
	mi := &file_public_input_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeBRequest.ProtoReflect.Descriptor instead.
func (*NodeBRequest) Descriptor() ([]byte, []int) {
	return file_public_input_proto_rawDescGZIP(), []int{3}
}

func (x *NodeBRequest) GetGreeting() string {
	if x != nil {
		return x.Greeting
	}
	return ""
}

func (x *NodeBRequest) GetTesting() *Testing {
	if x != nil {
		return x.Testing
	}
	return nil
}

func (x *NodeBRequest) GetTestingResults() *TestingResults {
	if x != nil {
		return x.TestingResults
	}
	return nil
}

type NodeCRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Greeting       string          `protobuf:"bytes,1,opt,name=greeting,proto3" json:"greeting,omitempty"`
	Testing        *Testing        `protobuf:"bytes,2,opt,name=testing,proto3" json:"testing,omitempty"`
	TestingResults *TestingResults `protobuf:"bytes,3,opt,name=testing_results,json=testingResults,proto3" json:"testing_results,omitempty"`
}

func (x *NodeCRequest) Reset() {
	*x = NodeCRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_public_input_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeCRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeCRequest) ProtoMessage() {}

func (x *NodeCRequest) ProtoReflect() protoreflect.Message {
	mi := &file_public_input_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeCRequest.ProtoReflect.Descriptor instead.
func (*NodeCRequest) Descriptor() ([]byte, []int) {
	return file_public_input_proto_rawDescGZIP(), []int{4}
}

func (x *NodeCRequest) GetGreeting() string {
	if x != nil {
		return x.Greeting
	}
	return ""
}

func (x *NodeCRequest) GetTesting() *Testing {
	if x != nil {
		return x.Testing
	}
	return nil
}

func (x *NodeCRequest) GetTestingResults() *TestingResults {
	if x != nil {
		return x.TestingResults
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Greeting       string          `protobuf:"bytes,1,opt,name=greeting,proto3" json:"greeting,omitempty"`
	Testing        *Testing        `protobuf:"bytes,2,opt,name=testing,proto3" json:"testing,omitempty"`
	TestingResults *TestingResults `protobuf:"bytes,3,opt,name=testing_results,json=testingResults,proto3" json:"testing_results,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_public_input_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_public_input_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_public_input_proto_rawDescGZIP(), []int{5}
}

func (x *Response) GetGreeting() string {
	if x != nil {
		return x.Greeting
	}
	return ""
}

func (x *Response) GetTesting() *Testing {
	if x != nil {
		return x.Testing
	}
	return nil
}

func (x *Response) GetTestingResults() *TestingResults {
	if x != nil {
		return x.TestingResults
	}
	return nil
}

var File_public_input_proto protoreflect.FileDescriptor

var file_public_input_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x74, 0x0a, 0x07, 0x54, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x24, 0x0a, 0x0e, 0x69, 0x73, 0x5f, 0x65, 0x61, 0x72, 0x6c,
	0x79, 0x5f, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x69,
	0x73, 0x45, 0x61, 0x72, 0x6c, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x22, 0x0a, 0x0d, 0x69,
	0x73, 0x5f, 0x65, 0x61, 0x72, 0x6c, 0x79, 0x5f, 0x65, 0x78, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0b, 0x69, 0x73, 0x45, 0x61, 0x72, 0x6c, 0x79, 0x45, 0x78, 0x69, 0x74, 0x12,
	0x1f, 0x0a, 0x0b, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x73,
	0x22, 0x40, 0x0a, 0x0e, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x73, 0x12, 0x2e, 0x0a, 0x13, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x73, 0x5f, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x11, 0x74, 0x65, 0x73, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x73, 0x53, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x22, 0x46, 0x0a, 0x07, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x27, 0x0a, 0x07, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6e,
	0x67, 0x52, 0x07, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x22, 0x92, 0x01, 0x0a, 0x0c, 0x4e,
	0x6f, 0x64, 0x65, 0x42, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x67,
	0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x67,
	0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x27, 0x0a, 0x07, 0x74, 0x65, 0x73, 0x74, 0x69,
	0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e,
	0x54, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x07, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67,
	0x12, 0x3d, 0x0a, 0x0f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x61, 0x69, 0x6e,
	0x2e, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x52,
	0x0e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x22,
	0x92, 0x01, 0x0a, 0x0c, 0x4e, 0x6f, 0x64, 0x65, 0x43, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x67, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x67, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x27, 0x0a, 0x07,
	0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x07, 0x74, 0x65,
	0x73, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x3d, 0x0a, 0x0f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67,
	0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x52, 0x0e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x22, 0x8e, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x67, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x67, 0x72, 0x65, 0x65, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x27, 0x0a,
	0x07, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x07, 0x74,
	0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x3d, 0x0a, 0x0f, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e,
	0x67, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x14, 0x2e, 0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x52, 0x0e, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x73, 0x32, 0x36, 0x0a, 0x0a, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x12, 0x28, 0x0a, 0x05, 0x47, 0x72, 0x65, 0x65, 0x74, 0x12, 0x0d, 0x2e, 0x6d,
	0x61, 0x69, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x6d, 0x61,
	0x69, 0x6e, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x09, 0x5a,
	0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_public_input_proto_rawDescOnce sync.Once
	file_public_input_proto_rawDescData = file_public_input_proto_rawDesc
)

func file_public_input_proto_rawDescGZIP() []byte {
	file_public_input_proto_rawDescOnce.Do(func() {
		file_public_input_proto_rawDescData = protoimpl.X.CompressGZIP(file_public_input_proto_rawDescData)
	})
	return file_public_input_proto_rawDescData
}

var file_public_input_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_public_input_proto_goTypes = []interface{}{
	(*Testing)(nil),        // 0: main.Testing
	(*TestingResults)(nil), // 1: main.TestingResults
	(*Request)(nil),        // 2: main.Request
	(*NodeBRequest)(nil),   // 3: main.NodeBRequest
	(*NodeCRequest)(nil),   // 4: main.NodeCRequest
	(*Response)(nil),       // 5: main.Response
}
var file_public_input_proto_depIdxs = []int32{
	0, // 0: main.Request.testing:type_name -> main.Testing
	0, // 1: main.NodeBRequest.testing:type_name -> main.Testing
	1, // 2: main.NodeBRequest.testing_results:type_name -> main.TestingResults
	0, // 3: main.NodeCRequest.testing:type_name -> main.Testing
	1, // 4: main.NodeCRequest.testing_results:type_name -> main.TestingResults
	0, // 5: main.Response.testing:type_name -> main.Testing
	1, // 6: main.Response.testing_results:type_name -> main.TestingResults
	2, // 7: main.Entrypoint.Greet:input_type -> main.Request
	5, // 8: main.Entrypoint.Greet:output_type -> main.Response
	8, // [8:9] is the sub-list for method output_type
	7, // [7:8] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_public_input_proto_init() }
func file_public_input_proto_init() {
	if File_public_input_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_public_input_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Testing); i {
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
		file_public_input_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestingResults); i {
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
		file_public_input_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_public_input_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeBRequest); i {
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
		file_public_input_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeCRequest); i {
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
		file_public_input_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
			RawDescriptor: file_public_input_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_public_input_proto_goTypes,
		DependencyIndexes: file_public_input_proto_depIdxs,
		MessageInfos:      file_public_input_proto_msgTypes,
	}.Build()
	File_public_input_proto = out.File
	file_public_input_proto_rawDesc = nil
	file_public_input_proto_goTypes = nil
	file_public_input_proto_depIdxs = nil
}

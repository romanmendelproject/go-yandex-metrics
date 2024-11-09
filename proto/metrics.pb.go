// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.12.4
// source: proto/metrics.proto

package go_yandex_metrics

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

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID    string  `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`          // имя метрики
	MType string  `protobuf:"bytes,2,opt,name=MType,proto3" json:"MType,omitempty"`    // параметр, принимающий значение gauge или counter
	Delta int64   `protobuf:"zigzag64,3,opt,name=Delta,proto3" json:"Delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `protobuf:"fixed64,4,opt,name=Value,proto3" json:"Value,omitempty"`  // значение метрики в случае передачи gauge
}

func (x *Metric) Reset() {
	*x = Metric{}
	mi := &file_proto_metrics_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Metric) GetMType() string {
	if x != nil {
		return x.MType
	}
	return ""
}

func (x *Metric) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Metric) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type ValueGaugeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *ValueGaugeRequest) Reset() {
	*x = ValueGaugeRequest{}
	mi := &file_proto_metrics_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValueGaugeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValueGaugeRequest) ProtoMessage() {}

func (x *ValueGaugeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValueGaugeRequest.ProtoReflect.Descriptor instead.
func (*ValueGaugeRequest) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{1}
}

func (x *ValueGaugeRequest) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

type ValueGaugeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value float64 `protobuf:"fixed64,1,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *ValueGaugeResponse) Reset() {
	*x = ValueGaugeResponse{}
	mi := &file_proto_metrics_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValueGaugeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValueGaugeResponse) ProtoMessage() {}

func (x *ValueGaugeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValueGaugeResponse.ProtoReflect.Descriptor instead.
func (*ValueGaugeResponse) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{2}
}

func (x *ValueGaugeResponse) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type ValueCounterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *ValueCounterRequest) Reset() {
	*x = ValueCounterRequest{}
	mi := &file_proto_metrics_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValueCounterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValueCounterRequest) ProtoMessage() {}

func (x *ValueCounterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValueCounterRequest.ProtoReflect.Descriptor instead.
func (*ValueCounterRequest) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{3}
}

func (x *ValueCounterRequest) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

type ValueCounterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Delta int64 `protobuf:"zigzag64,1,opt,name=Delta,proto3" json:"Delta,omitempty"`
}

func (x *ValueCounterResponse) Reset() {
	*x = ValueCounterResponse{}
	mi := &file_proto_metrics_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValueCounterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValueCounterResponse) ProtoMessage() {}

func (x *ValueCounterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValueCounterResponse.ProtoReflect.Descriptor instead.
func (*ValueCounterResponse) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{4}
}

func (x *ValueCounterResponse) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

type UpdateBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric []*Metric `protobuf:"bytes,1,rep,name=metric,proto3" json:"metric,omitempty"`
}

func (x *UpdateBatchRequest) Reset() {
	*x = UpdateBatchRequest{}
	mi := &file_proto_metrics_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBatchRequest) ProtoMessage() {}

func (x *UpdateBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBatchRequest.ProtoReflect.Descriptor instead.
func (*UpdateBatchRequest) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateBatchRequest) GetMetric() []*Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

type UpdateBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metric []*Metric `protobuf:"bytes,1,rep,name=metric,proto3" json:"metric,omitempty"`
}

func (x *UpdateBatchResponse) Reset() {
	*x = UpdateBatchResponse{}
	mi := &file_proto_metrics_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateBatchResponse) ProtoMessage() {}

func (x *UpdateBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrics_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateBatchResponse.ProtoReflect.Descriptor instead.
func (*UpdateBatchResponse) Descriptor() ([]byte, []int) {
	return file_proto_metrics_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateBatchResponse) GetMetric() []*Metric {
	if x != nil {
		return x.Metric
	}
	return nil
}

var File_proto_metrics_proto protoreflect.FileDescriptor

var file_proto_metrics_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x64, 0x65, 0x6d, 0x6f, 0x22, 0x5a, 0x0a, 0x06, 0x4d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x4d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x44,
	0x65, 0x6c, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x12, 0x52, 0x05, 0x44, 0x65, 0x6c, 0x74,
	0x61, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x23, 0x0a, 0x11, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x47, 0x61, 0x75, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x22, 0x2a, 0x0a, 0x12,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x47, 0x61, 0x75, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x25, 0x0a, 0x13, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x22,
	0x2c, 0x0a, 0x14, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x44, 0x65, 0x6c, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x12, 0x52, 0x05, 0x44, 0x65, 0x6c, 0x74, 0x61, 0x22, 0x3a, 0x0a,
	0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x22, 0x3b, 0x0a, 0x13, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x24, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0c, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x52, 0x06,
	0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x32, 0xd5, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x12, 0x3f, 0x0a, 0x0a, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x47, 0x61, 0x75, 0x67, 0x65,
	0x12, 0x17, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x47, 0x61, 0x75,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x64, 0x65, 0x6d, 0x6f,
	0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x47, 0x61, 0x75, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x45, 0x0a, 0x0c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x65, 0x72, 0x12, 0x19, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x0b, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x18, 0x2e, 0x64, 0x65, 0x6d, 0x6f,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x64, 0x65, 0x6d, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x31,
	0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x6d,
	0x61, 0x6e, 0x6d, 0x65, 0x6e, 0x64, 0x65, 0x6c, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f,
	0x67, 0x6f, 0x2d, 0x79, 0x61, 0x6e, 0x64, 0x65, 0x78, 0x2d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_metrics_proto_rawDescOnce sync.Once
	file_proto_metrics_proto_rawDescData = file_proto_metrics_proto_rawDesc
)

func file_proto_metrics_proto_rawDescGZIP() []byte {
	file_proto_metrics_proto_rawDescOnce.Do(func() {
		file_proto_metrics_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_metrics_proto_rawDescData)
	})
	return file_proto_metrics_proto_rawDescData
}

var file_proto_metrics_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_metrics_proto_goTypes = []any{
	(*Metric)(nil),               // 0: demo.Metric
	(*ValueGaugeRequest)(nil),    // 1: demo.ValueGaugeRequest
	(*ValueGaugeResponse)(nil),   // 2: demo.ValueGaugeResponse
	(*ValueCounterRequest)(nil),  // 3: demo.ValueCounterRequest
	(*ValueCounterResponse)(nil), // 4: demo.ValueCounterResponse
	(*UpdateBatchRequest)(nil),   // 5: demo.UpdateBatchRequest
	(*UpdateBatchResponse)(nil),  // 6: demo.UpdateBatchResponse
}
var file_proto_metrics_proto_depIdxs = []int32{
	0, // 0: demo.UpdateBatchRequest.metric:type_name -> demo.Metric
	0, // 1: demo.UpdateBatchResponse.metric:type_name -> demo.Metric
	1, // 2: demo.Metrics.ValueGauge:input_type -> demo.ValueGaugeRequest
	3, // 3: demo.Metrics.ValueCounter:input_type -> demo.ValueCounterRequest
	5, // 4: demo.Metrics.UpdateBatch:input_type -> demo.UpdateBatchRequest
	2, // 5: demo.Metrics.ValueGauge:output_type -> demo.ValueGaugeResponse
	4, // 6: demo.Metrics.ValueCounter:output_type -> demo.ValueCounterResponse
	6, // 7: demo.Metrics.UpdateBatch:output_type -> demo.UpdateBatchResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_metrics_proto_init() }
func file_proto_metrics_proto_init() {
	if File_proto_metrics_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_metrics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_metrics_proto_goTypes,
		DependencyIndexes: file_proto_metrics_proto_depIdxs,
		MessageInfos:      file_proto_metrics_proto_msgTypes,
	}.Build()
	File_proto_metrics_proto = out.File
	file_proto_metrics_proto_rawDesc = nil
	file_proto_metrics_proto_goTypes = nil
	file_proto_metrics_proto_depIdxs = nil
}

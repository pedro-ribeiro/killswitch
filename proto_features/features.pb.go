// Code generated by protoc-gen-go. DO NOT EDIT.
// source: features.proto

package features

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Feature struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Feature) Reset()         { *m = Feature{} }
func (m *Feature) String() string { return proto.CompactTextString(m) }
func (*Feature) ProtoMessage()    {}
func (*Feature) Descriptor() ([]byte, []int) {
	return fileDescriptor_2216f05915163cdf, []int{0}
}

func (m *Feature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Feature.Unmarshal(m, b)
}
func (m *Feature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Feature.Marshal(b, m, deterministic)
}
func (m *Feature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Feature.Merge(m, src)
}
func (m *Feature) XXX_Size() int {
	return xxx_messageInfo_Feature.Size(m)
}
func (m *Feature) XXX_DiscardUnknown() {
	xxx_messageInfo_Feature.DiscardUnknown(m)
}

var xxx_messageInfo_Feature proto.InternalMessageInfo

type FeaturesRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FeaturesRequest) Reset()         { *m = FeaturesRequest{} }
func (m *FeaturesRequest) String() string { return proto.CompactTextString(m) }
func (*FeaturesRequest) ProtoMessage()    {}
func (*FeaturesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2216f05915163cdf, []int{1}
}

func (m *FeaturesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeaturesRequest.Unmarshal(m, b)
}
func (m *FeaturesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeaturesRequest.Marshal(b, m, deterministic)
}
func (m *FeaturesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeaturesRequest.Merge(m, src)
}
func (m *FeaturesRequest) XXX_Size() int {
	return xxx_messageInfo_FeaturesRequest.Size(m)
}
func (m *FeaturesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FeaturesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FeaturesRequest proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Feature)(nil), "Feature")
	proto.RegisterType((*FeaturesRequest)(nil), "FeaturesRequest")
}

func init() { proto.RegisterFile("features.proto", fileDescriptor_2216f05915163cdf) }

var fileDescriptor_2216f05915163cdf = []byte{
	// 87 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0x4b, 0x4d, 0x2c,
	0x29, 0x2d, 0x4a, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xe2, 0xe4, 0x62, 0x77, 0x83,
	0x88, 0x28, 0x09, 0x72, 0xf1, 0x43, 0x99, 0xc5, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x46,
	0x96, 0x5c, 0x1c, 0x30, 0x21, 0x21, 0x5d, 0x2e, 0x6e, 0xf7, 0xd4, 0x12, 0x38, 0x57, 0x40, 0x0f,
	0x4d, 0xb1, 0x14, 0x07, 0x4c, 0x44, 0x89, 0xc1, 0x80, 0x11, 0x10, 0x00, 0x00, 0xff, 0xff, 0x97,
	0x0f, 0x92, 0x3f, 0x69, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FeaturesClient is the client API for Features service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FeaturesClient interface {
	GetFeatures(ctx context.Context, in *FeaturesRequest, opts ...grpc.CallOption) (Features_GetFeaturesClient, error)
}

type featuresClient struct {
	cc *grpc.ClientConn
}

func NewFeaturesClient(cc *grpc.ClientConn) FeaturesClient {
	return &featuresClient{cc}
}

func (c *featuresClient) GetFeatures(ctx context.Context, in *FeaturesRequest, opts ...grpc.CallOption) (Features_GetFeaturesClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Features_serviceDesc.Streams[0], "/Features/GetFeatures", opts...)
	if err != nil {
		return nil, err
	}
	x := &featuresGetFeaturesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Features_GetFeaturesClient interface {
	Recv() (*Feature, error)
	grpc.ClientStream
}

type featuresGetFeaturesClient struct {
	grpc.ClientStream
}

func (x *featuresGetFeaturesClient) Recv() (*Feature, error) {
	m := new(Feature)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FeaturesServer is the server API for Features service.
type FeaturesServer interface {
	GetFeatures(*FeaturesRequest, Features_GetFeaturesServer) error
}

// UnimplementedFeaturesServer can be embedded to have forward compatible implementations.
type UnimplementedFeaturesServer struct {
}

func (*UnimplementedFeaturesServer) GetFeatures(req *FeaturesRequest, srv Features_GetFeaturesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetFeatures not implemented")
}

func RegisterFeaturesServer(s *grpc.Server, srv FeaturesServer) {
	s.RegisterService(&_Features_serviceDesc, srv)
}

func _Features_GetFeatures_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FeaturesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FeaturesServer).GetFeatures(m, &featuresGetFeaturesServer{stream})
}

type Features_GetFeaturesServer interface {
	Send(*Feature) error
	grpc.ServerStream
}

type featuresGetFeaturesServer struct {
	grpc.ServerStream
}

func (x *featuresGetFeaturesServer) Send(m *Feature) error {
	return x.ServerStream.SendMsg(m)
}

var _Features_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Features",
	HandlerType: (*FeaturesServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetFeatures",
			Handler:       _Features_GetFeatures_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "features.proto",
}

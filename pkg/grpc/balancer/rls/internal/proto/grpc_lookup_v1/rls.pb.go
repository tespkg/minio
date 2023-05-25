// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc/rls/grpc_lookup_v1/rls.proto

package grpc_lookup_v1

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "github.com/minio/minio/pkg/grpc"
	codes "github.com/minio/minio/pkg/grpc/codes"
	status "github.com/minio/minio/pkg/grpc/status"
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

type RouteLookupRequest struct {
	// Full host name of the target server, e.g. firestore.googleapis.com.
	// Only set for gRPC requests; HTTP requests must use key_map explicitly.
	Server string `protobuf:"bytes,1,opt,name=server,proto3" json:"server,omitempty"`
	// Full path of the request, i.e. "/service/method".
	// Only set for gRPC requests; HTTP requests must use key_map explicitly.
	Path string `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	// Target type allows the client to specify what kind of target format it
	// would like from RLS to allow it to find the regional server, e.g. "grpc".
	TargetType string `protobuf:"bytes,3,opt,name=target_type,json=targetType,proto3" json:"target_type,omitempty"`
	// Map of key values extracted via key builders for the gRPC or HTTP request.
	KeyMap               map[string]string `protobuf:"bytes,4,rep,name=key_map,json=keyMap,proto3" json:"key_map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *RouteLookupRequest) Reset()         { *m = RouteLookupRequest{} }
func (m *RouteLookupRequest) String() string { return proto.CompactTextString(m) }
func (*RouteLookupRequest) ProtoMessage()    {}
func (*RouteLookupRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5fe9649e373b9d12, []int{0}
}

func (m *RouteLookupRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RouteLookupRequest.Unmarshal(m, b)
}
func (m *RouteLookupRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RouteLookupRequest.Marshal(b, m, deterministic)
}
func (m *RouteLookupRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RouteLookupRequest.Merge(m, src)
}
func (m *RouteLookupRequest) XXX_Size() int {
	return xxx_messageInfo_RouteLookupRequest.Size(m)
}
func (m *RouteLookupRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RouteLookupRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RouteLookupRequest proto.InternalMessageInfo

func (m *RouteLookupRequest) GetServer() string {
	if m != nil {
		return m.Server
	}
	return ""
}

func (m *RouteLookupRequest) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *RouteLookupRequest) GetTargetType() string {
	if m != nil {
		return m.TargetType
	}
	return ""
}

func (m *RouteLookupRequest) GetKeyMap() map[string]string {
	if m != nil {
		return m.KeyMap
	}
	return nil
}

type RouteLookupResponse struct {
	// Actual addressable entity to use for routing decision, using syntax
	// requested by the request target_type.
	Target string `protobuf:"bytes,1,opt,name=target,proto3" json:"target,omitempty"`
	// Optional header value to pass along to AFE in the X-Google-RLS-Data header.
	// Cached with "target" and sent with all requests that match the request key.
	// Allows the RLS to pass its work product to the eventual target.
	HeaderData           string   `protobuf:"bytes,2,opt,name=header_data,json=headerData,proto3" json:"header_data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RouteLookupResponse) Reset()         { *m = RouteLookupResponse{} }
func (m *RouteLookupResponse) String() string { return proto.CompactTextString(m) }
func (*RouteLookupResponse) ProtoMessage()    {}
func (*RouteLookupResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5fe9649e373b9d12, []int{1}
}

func (m *RouteLookupResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RouteLookupResponse.Unmarshal(m, b)
}
func (m *RouteLookupResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RouteLookupResponse.Marshal(b, m, deterministic)
}
func (m *RouteLookupResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RouteLookupResponse.Merge(m, src)
}
func (m *RouteLookupResponse) XXX_Size() int {
	return xxx_messageInfo_RouteLookupResponse.Size(m)
}
func (m *RouteLookupResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RouteLookupResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RouteLookupResponse proto.InternalMessageInfo

func (m *RouteLookupResponse) GetTarget() string {
	if m != nil {
		return m.Target
	}
	return ""
}

func (m *RouteLookupResponse) GetHeaderData() string {
	if m != nil {
		return m.HeaderData
	}
	return ""
}

func init() {
	proto.RegisterType((*RouteLookupRequest)(nil), "grpc.lookup.v1.RouteLookupRequest")
	proto.RegisterMapType((map[string]string)(nil), "grpc.lookup.v1.RouteLookupRequest.KeyMapEntry")
	proto.RegisterType((*RouteLookupResponse)(nil), "grpc.lookup.v1.RouteLookupResponse")
}

func init() { proto.RegisterFile("grpc/rls/grpc_lookup_v1/rls.proto", fileDescriptor_5fe9649e373b9d12) }

var fileDescriptor_5fe9649e373b9d12 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xc1, 0x4b, 0xc3, 0x30,
	0x14, 0xc6, 0xed, 0x36, 0xa7, 0xbe, 0x82, 0x68, 0x14, 0x29, 0xbb, 0x38, 0xeb, 0x65, 0x07, 0xc9,
	0xd8, 0xbc, 0xa8, 0xc7, 0xa1, 0x78, 0xd0, 0xc9, 0xa8, 0x1e, 0xc4, 0x4b, 0x89, 0xdb, 0x23, 0x1b,
	0xad, 0x4d, 0x4c, 0xd3, 0x42, 0xff, 0x60, 0xff, 0x0f, 0x49, 0x52, 0x61, 0x9d, 0xa0, 0xb7, 0xf7,
	0xfd, 0xde, 0x23, 0xf9, 0xbe, 0xe4, 0xc1, 0x19, 0x57, 0x72, 0x3e, 0x54, 0x69, 0x3e, 0x34, 0x45,
	0x9c, 0x0a, 0x91, 0x14, 0x32, 0x2e, 0x47, 0x06, 0x51, 0xa9, 0x84, 0x16, 0x64, 0xdf, 0x74, 0xa8,
	0xeb, 0xd0, 0x72, 0x14, 0x7e, 0x79, 0x40, 0x22, 0x51, 0x68, 0x7c, 0xb4, 0x28, 0xc2, 0xcf, 0x02,
	0x73, 0x4d, 0x4e, 0xa0, 0x9b, 0xa3, 0x2a, 0x51, 0x05, 0x5e, 0xdf, 0x1b, 0xec, 0x45, 0xb5, 0x22,
	0x04, 0x3a, 0x92, 0xe9, 0x65, 0xd0, 0xb2, 0xd4, 0xd6, 0xe4, 0x14, 0x7c, 0xcd, 0x14, 0x47, 0x1d,
	0xeb, 0x4a, 0x62, 0xd0, 0xb6, 0x2d, 0x70, 0xe8, 0xa5, 0x92, 0x48, 0xee, 0x61, 0x27, 0xc1, 0x2a,
	0xfe, 0x60, 0x32, 0xe8, 0xf4, 0xdb, 0x03, 0x7f, 0x4c, 0x69, 0xd3, 0x05, 0xfd, 0xed, 0x80, 0x3e,
	0x60, 0x35, 0x65, 0xf2, 0x2e, 0xd3, 0xaa, 0x8a, 0xba, 0x89, 0x15, 0xbd, 0x6b, 0xf0, 0xd7, 0x30,
	0x39, 0x80, 0x76, 0x82, 0x55, 0xed, 0xd0, 0x94, 0xe4, 0x18, 0xb6, 0x4b, 0x96, 0x16, 0x58, 0xfb,
	0x73, 0xe2, 0xa6, 0x75, 0xe5, 0x85, 0x4f, 0x70, 0xd4, 0xb8, 0x24, 0x97, 0x22, 0xcb, 0xd1, 0xe4,
	0x74, 0x46, 0x7f, 0x72, 0x3a, 0x65, 0x32, 0x2d, 0x91, 0x2d, 0x50, 0xc5, 0x0b, 0xa6, 0x59, 0x7d,
	0x1c, 0x38, 0x74, 0xcb, 0x34, 0x1b, 0x67, 0x8d, 0x67, 0x7b, 0x46, 0x55, 0xae, 0xe6, 0x48, 0x5e,
	0xc1, 0x5f, 0xa3, 0x24, 0xfc, 0x3f, 0x67, 0xef, 0xfc, 0xcf, 0x19, 0x67, 0x33, 0xdc, 0x9a, 0x4c,
	0xe1, 0x70, 0x25, 0x36, 0x46, 0x27, 0xbb, 0x51, 0x9a, 0xcf, 0xcc, 0xb7, 0xce, 0xbc, 0xb7, 0x0b,
	0x2e, 0x04, 0x4f, 0x91, 0x72, 0x91, 0xb2, 0x8c, 0x53, 0xa1, 0xb8, 0x5d, 0x82, 0xa1, 0x9b, 0xde,
	0x58, 0x88, 0xf7, 0xae, 0xdd, 0x86, 0xcb, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xba, 0x10, 0x2d,
	0xb5, 0x32, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RouteLookupServiceClient is the client API for RouteLookupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/github.com/minio/minio/pkg/grpc#ClientConn.NewStream.
type RouteLookupServiceClient interface {
	// Lookup returns a target for a single key.
	RouteLookup(ctx context.Context, in *RouteLookupRequest, opts ...grpc.CallOption) (*RouteLookupResponse, error)
}

type routeLookupServiceClient struct {
	cc *grpc.ClientConn
}

func NewRouteLookupServiceClient(cc *grpc.ClientConn) RouteLookupServiceClient {
	return &routeLookupServiceClient{cc}
}

func (c *routeLookupServiceClient) RouteLookup(ctx context.Context, in *RouteLookupRequest, opts ...grpc.CallOption) (*RouteLookupResponse, error) {
	out := new(RouteLookupResponse)
	err := c.cc.Invoke(ctx, "/grpc.lookup.v1.RouteLookupService/RouteLookup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RouteLookupServiceServer is the server API for RouteLookupService service.
type RouteLookupServiceServer interface {
	// Lookup returns a target for a single key.
	RouteLookup(context.Context, *RouteLookupRequest) (*RouteLookupResponse, error)
}

// UnimplementedRouteLookupServiceServer can be embedded to have forward compatible implementations.
type UnimplementedRouteLookupServiceServer struct {
}

func (*UnimplementedRouteLookupServiceServer) RouteLookup(ctx context.Context, req *RouteLookupRequest) (*RouteLookupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RouteLookup not implemented")
}

func RegisterRouteLookupServiceServer(s *grpc.Server, srv RouteLookupServiceServer) {
	s.RegisterService(&_RouteLookupService_serviceDesc, srv)
}

func _RouteLookupService_RouteLookup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RouteLookupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouteLookupServiceServer).RouteLookup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.lookup.v1.RouteLookupService/RouteLookup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouteLookupServiceServer).RouteLookup(ctx, req.(*RouteLookupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RouteLookupService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.lookup.v1.RouteLookupService",
	HandlerType: (*RouteLookupServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RouteLookup",
			Handler:    _RouteLookupService_RouteLookup_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/rls/grpc_lookup_v1/rls.proto",
}

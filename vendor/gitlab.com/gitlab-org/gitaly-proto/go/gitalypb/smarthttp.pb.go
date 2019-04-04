// Code generated by protoc-gen-go. DO NOT EDIT.
// source: smarthttp.proto

package gitalypb // import "gitlab.com/gitlab-org/gitaly-proto/go/gitalypb"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type InfoRefsRequest struct {
	Repository *Repository `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	// Parameters to use with git -c (key=value pairs)
	GitConfigOptions []string `protobuf:"bytes,2,rep,name=git_config_options,json=gitConfigOptions,proto3" json:"git_config_options,omitempty"`
	// Git protocol version
	GitProtocol          string   `protobuf:"bytes,3,opt,name=git_protocol,json=gitProtocol,proto3" json:"git_protocol,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoRefsRequest) Reset()         { *m = InfoRefsRequest{} }
func (m *InfoRefsRequest) String() string { return proto.CompactTextString(m) }
func (*InfoRefsRequest) ProtoMessage()    {}
func (*InfoRefsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_smarthttp_d15d08ac1e07ff5f, []int{0}
}
func (m *InfoRefsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoRefsRequest.Unmarshal(m, b)
}
func (m *InfoRefsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoRefsRequest.Marshal(b, m, deterministic)
}
func (dst *InfoRefsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoRefsRequest.Merge(dst, src)
}
func (m *InfoRefsRequest) XXX_Size() int {
	return xxx_messageInfo_InfoRefsRequest.Size(m)
}
func (m *InfoRefsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoRefsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InfoRefsRequest proto.InternalMessageInfo

func (m *InfoRefsRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *InfoRefsRequest) GetGitConfigOptions() []string {
	if m != nil {
		return m.GitConfigOptions
	}
	return nil
}

func (m *InfoRefsRequest) GetGitProtocol() string {
	if m != nil {
		return m.GitProtocol
	}
	return ""
}

type InfoRefsResponse struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoRefsResponse) Reset()         { *m = InfoRefsResponse{} }
func (m *InfoRefsResponse) String() string { return proto.CompactTextString(m) }
func (*InfoRefsResponse) ProtoMessage()    {}
func (*InfoRefsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_smarthttp_d15d08ac1e07ff5f, []int{1}
}
func (m *InfoRefsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoRefsResponse.Unmarshal(m, b)
}
func (m *InfoRefsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoRefsResponse.Marshal(b, m, deterministic)
}
func (dst *InfoRefsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoRefsResponse.Merge(dst, src)
}
func (m *InfoRefsResponse) XXX_Size() int {
	return xxx_messageInfo_InfoRefsResponse.Size(m)
}
func (m *InfoRefsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoRefsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InfoRefsResponse proto.InternalMessageInfo

func (m *InfoRefsResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type PostUploadPackRequest struct {
	// repository should only be present in the first message of the stream
	Repository *Repository `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	// Raw data to be copied to stdin of 'git upload-pack'
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// Parameters to use with git -c (key=value pairs)
	GitConfigOptions []string `protobuf:"bytes,3,rep,name=git_config_options,json=gitConfigOptions,proto3" json:"git_config_options,omitempty"`
	// Git protocol version
	GitProtocol          string   `protobuf:"bytes,4,opt,name=git_protocol,json=gitProtocol,proto3" json:"git_protocol,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PostUploadPackRequest) Reset()         { *m = PostUploadPackRequest{} }
func (m *PostUploadPackRequest) String() string { return proto.CompactTextString(m) }
func (*PostUploadPackRequest) ProtoMessage()    {}
func (*PostUploadPackRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_smarthttp_d15d08ac1e07ff5f, []int{2}
}
func (m *PostUploadPackRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PostUploadPackRequest.Unmarshal(m, b)
}
func (m *PostUploadPackRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PostUploadPackRequest.Marshal(b, m, deterministic)
}
func (dst *PostUploadPackRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PostUploadPackRequest.Merge(dst, src)
}
func (m *PostUploadPackRequest) XXX_Size() int {
	return xxx_messageInfo_PostUploadPackRequest.Size(m)
}
func (m *PostUploadPackRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PostUploadPackRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PostUploadPackRequest proto.InternalMessageInfo

func (m *PostUploadPackRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *PostUploadPackRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *PostUploadPackRequest) GetGitConfigOptions() []string {
	if m != nil {
		return m.GitConfigOptions
	}
	return nil
}

func (m *PostUploadPackRequest) GetGitProtocol() string {
	if m != nil {
		return m.GitProtocol
	}
	return ""
}

type PostUploadPackResponse struct {
	// Raw data from stdout of 'git upload-pack'
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PostUploadPackResponse) Reset()         { *m = PostUploadPackResponse{} }
func (m *PostUploadPackResponse) String() string { return proto.CompactTextString(m) }
func (*PostUploadPackResponse) ProtoMessage()    {}
func (*PostUploadPackResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_smarthttp_d15d08ac1e07ff5f, []int{3}
}
func (m *PostUploadPackResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PostUploadPackResponse.Unmarshal(m, b)
}
func (m *PostUploadPackResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PostUploadPackResponse.Marshal(b, m, deterministic)
}
func (dst *PostUploadPackResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PostUploadPackResponse.Merge(dst, src)
}
func (m *PostUploadPackResponse) XXX_Size() int {
	return xxx_messageInfo_PostUploadPackResponse.Size(m)
}
func (m *PostUploadPackResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PostUploadPackResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PostUploadPackResponse proto.InternalMessageInfo

func (m *PostUploadPackResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type PostReceivePackRequest struct {
	// repository should only be present in the first message of the stream
	Repository *Repository `protobuf:"bytes,1,opt,name=repository,proto3" json:"repository,omitempty"`
	// Raw data to be copied to stdin of 'git receive-pack'
	Data []byte `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// gl_id, gl_repository, and gl_username become env variables, used by the Git {pre,post}-receive
	// hooks. They should only be present in the first message of the stream.
	GlId         string `protobuf:"bytes,3,opt,name=gl_id,json=glId,proto3" json:"gl_id,omitempty"`
	GlRepository string `protobuf:"bytes,4,opt,name=gl_repository,json=glRepository,proto3" json:"gl_repository,omitempty"`
	GlUsername   string `protobuf:"bytes,5,opt,name=gl_username,json=glUsername,proto3" json:"gl_username,omitempty"`
	// Git protocol version
	GitProtocol string `protobuf:"bytes,6,opt,name=git_protocol,json=gitProtocol,proto3" json:"git_protocol,omitempty"`
	// Parameters to use with git -c (key=value pairs)
	GitConfigOptions     []string `protobuf:"bytes,7,rep,name=git_config_options,json=gitConfigOptions,proto3" json:"git_config_options,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PostReceivePackRequest) Reset()         { *m = PostReceivePackRequest{} }
func (m *PostReceivePackRequest) String() string { return proto.CompactTextString(m) }
func (*PostReceivePackRequest) ProtoMessage()    {}
func (*PostReceivePackRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_smarthttp_d15d08ac1e07ff5f, []int{4}
}
func (m *PostReceivePackRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PostReceivePackRequest.Unmarshal(m, b)
}
func (m *PostReceivePackRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PostReceivePackRequest.Marshal(b, m, deterministic)
}
func (dst *PostReceivePackRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PostReceivePackRequest.Merge(dst, src)
}
func (m *PostReceivePackRequest) XXX_Size() int {
	return xxx_messageInfo_PostReceivePackRequest.Size(m)
}
func (m *PostReceivePackRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PostReceivePackRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PostReceivePackRequest proto.InternalMessageInfo

func (m *PostReceivePackRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *PostReceivePackRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *PostReceivePackRequest) GetGlId() string {
	if m != nil {
		return m.GlId
	}
	return ""
}

func (m *PostReceivePackRequest) GetGlRepository() string {
	if m != nil {
		return m.GlRepository
	}
	return ""
}

func (m *PostReceivePackRequest) GetGlUsername() string {
	if m != nil {
		return m.GlUsername
	}
	return ""
}

func (m *PostReceivePackRequest) GetGitProtocol() string {
	if m != nil {
		return m.GitProtocol
	}
	return ""
}

func (m *PostReceivePackRequest) GetGitConfigOptions() []string {
	if m != nil {
		return m.GitConfigOptions
	}
	return nil
}

type PostReceivePackResponse struct {
	// Raw data from stdout of 'git receive-pack'
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PostReceivePackResponse) Reset()         { *m = PostReceivePackResponse{} }
func (m *PostReceivePackResponse) String() string { return proto.CompactTextString(m) }
func (*PostReceivePackResponse) ProtoMessage()    {}
func (*PostReceivePackResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_smarthttp_d15d08ac1e07ff5f, []int{5}
}
func (m *PostReceivePackResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PostReceivePackResponse.Unmarshal(m, b)
}
func (m *PostReceivePackResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PostReceivePackResponse.Marshal(b, m, deterministic)
}
func (dst *PostReceivePackResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PostReceivePackResponse.Merge(dst, src)
}
func (m *PostReceivePackResponse) XXX_Size() int {
	return xxx_messageInfo_PostReceivePackResponse.Size(m)
}
func (m *PostReceivePackResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PostReceivePackResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PostReceivePackResponse proto.InternalMessageInfo

func (m *PostReceivePackResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*InfoRefsRequest)(nil), "gitaly.InfoRefsRequest")
	proto.RegisterType((*InfoRefsResponse)(nil), "gitaly.InfoRefsResponse")
	proto.RegisterType((*PostUploadPackRequest)(nil), "gitaly.PostUploadPackRequest")
	proto.RegisterType((*PostUploadPackResponse)(nil), "gitaly.PostUploadPackResponse")
	proto.RegisterType((*PostReceivePackRequest)(nil), "gitaly.PostReceivePackRequest")
	proto.RegisterType((*PostReceivePackResponse)(nil), "gitaly.PostReceivePackResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SmartHTTPServiceClient is the client API for SmartHTTPService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SmartHTTPServiceClient interface {
	// The response body for GET /info/refs?service=git-upload-pack
	InfoRefsUploadPack(ctx context.Context, in *InfoRefsRequest, opts ...grpc.CallOption) (SmartHTTPService_InfoRefsUploadPackClient, error)
	// The response body for GET /info/refs?service=git-receive-pack
	InfoRefsReceivePack(ctx context.Context, in *InfoRefsRequest, opts ...grpc.CallOption) (SmartHTTPService_InfoRefsReceivePackClient, error)
	// Request and response body for POST /upload-pack
	PostUploadPack(ctx context.Context, opts ...grpc.CallOption) (SmartHTTPService_PostUploadPackClient, error)
	// Request and response body for POST /receive-pack
	PostReceivePack(ctx context.Context, opts ...grpc.CallOption) (SmartHTTPService_PostReceivePackClient, error)
}

type smartHTTPServiceClient struct {
	cc *grpc.ClientConn
}

func NewSmartHTTPServiceClient(cc *grpc.ClientConn) SmartHTTPServiceClient {
	return &smartHTTPServiceClient{cc}
}

func (c *smartHTTPServiceClient) InfoRefsUploadPack(ctx context.Context, in *InfoRefsRequest, opts ...grpc.CallOption) (SmartHTTPService_InfoRefsUploadPackClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SmartHTTPService_serviceDesc.Streams[0], "/gitaly.SmartHTTPService/InfoRefsUploadPack", opts...)
	if err != nil {
		return nil, err
	}
	x := &smartHTTPServiceInfoRefsUploadPackClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SmartHTTPService_InfoRefsUploadPackClient interface {
	Recv() (*InfoRefsResponse, error)
	grpc.ClientStream
}

type smartHTTPServiceInfoRefsUploadPackClient struct {
	grpc.ClientStream
}

func (x *smartHTTPServiceInfoRefsUploadPackClient) Recv() (*InfoRefsResponse, error) {
	m := new(InfoRefsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *smartHTTPServiceClient) InfoRefsReceivePack(ctx context.Context, in *InfoRefsRequest, opts ...grpc.CallOption) (SmartHTTPService_InfoRefsReceivePackClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SmartHTTPService_serviceDesc.Streams[1], "/gitaly.SmartHTTPService/InfoRefsReceivePack", opts...)
	if err != nil {
		return nil, err
	}
	x := &smartHTTPServiceInfoRefsReceivePackClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SmartHTTPService_InfoRefsReceivePackClient interface {
	Recv() (*InfoRefsResponse, error)
	grpc.ClientStream
}

type smartHTTPServiceInfoRefsReceivePackClient struct {
	grpc.ClientStream
}

func (x *smartHTTPServiceInfoRefsReceivePackClient) Recv() (*InfoRefsResponse, error) {
	m := new(InfoRefsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *smartHTTPServiceClient) PostUploadPack(ctx context.Context, opts ...grpc.CallOption) (SmartHTTPService_PostUploadPackClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SmartHTTPService_serviceDesc.Streams[2], "/gitaly.SmartHTTPService/PostUploadPack", opts...)
	if err != nil {
		return nil, err
	}
	x := &smartHTTPServicePostUploadPackClient{stream}
	return x, nil
}

type SmartHTTPService_PostUploadPackClient interface {
	Send(*PostUploadPackRequest) error
	Recv() (*PostUploadPackResponse, error)
	grpc.ClientStream
}

type smartHTTPServicePostUploadPackClient struct {
	grpc.ClientStream
}

func (x *smartHTTPServicePostUploadPackClient) Send(m *PostUploadPackRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *smartHTTPServicePostUploadPackClient) Recv() (*PostUploadPackResponse, error) {
	m := new(PostUploadPackResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *smartHTTPServiceClient) PostReceivePack(ctx context.Context, opts ...grpc.CallOption) (SmartHTTPService_PostReceivePackClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SmartHTTPService_serviceDesc.Streams[3], "/gitaly.SmartHTTPService/PostReceivePack", opts...)
	if err != nil {
		return nil, err
	}
	x := &smartHTTPServicePostReceivePackClient{stream}
	return x, nil
}

type SmartHTTPService_PostReceivePackClient interface {
	Send(*PostReceivePackRequest) error
	Recv() (*PostReceivePackResponse, error)
	grpc.ClientStream
}

type smartHTTPServicePostReceivePackClient struct {
	grpc.ClientStream
}

func (x *smartHTTPServicePostReceivePackClient) Send(m *PostReceivePackRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *smartHTTPServicePostReceivePackClient) Recv() (*PostReceivePackResponse, error) {
	m := new(PostReceivePackResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SmartHTTPServiceServer is the server API for SmartHTTPService service.
type SmartHTTPServiceServer interface {
	// The response body for GET /info/refs?service=git-upload-pack
	InfoRefsUploadPack(*InfoRefsRequest, SmartHTTPService_InfoRefsUploadPackServer) error
	// The response body for GET /info/refs?service=git-receive-pack
	InfoRefsReceivePack(*InfoRefsRequest, SmartHTTPService_InfoRefsReceivePackServer) error
	// Request and response body for POST /upload-pack
	PostUploadPack(SmartHTTPService_PostUploadPackServer) error
	// Request and response body for POST /receive-pack
	PostReceivePack(SmartHTTPService_PostReceivePackServer) error
}

func RegisterSmartHTTPServiceServer(s *grpc.Server, srv SmartHTTPServiceServer) {
	s.RegisterService(&_SmartHTTPService_serviceDesc, srv)
}

func _SmartHTTPService_InfoRefsUploadPack_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(InfoRefsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SmartHTTPServiceServer).InfoRefsUploadPack(m, &smartHTTPServiceInfoRefsUploadPackServer{stream})
}

type SmartHTTPService_InfoRefsUploadPackServer interface {
	Send(*InfoRefsResponse) error
	grpc.ServerStream
}

type smartHTTPServiceInfoRefsUploadPackServer struct {
	grpc.ServerStream
}

func (x *smartHTTPServiceInfoRefsUploadPackServer) Send(m *InfoRefsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _SmartHTTPService_InfoRefsReceivePack_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(InfoRefsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SmartHTTPServiceServer).InfoRefsReceivePack(m, &smartHTTPServiceInfoRefsReceivePackServer{stream})
}

type SmartHTTPService_InfoRefsReceivePackServer interface {
	Send(*InfoRefsResponse) error
	grpc.ServerStream
}

type smartHTTPServiceInfoRefsReceivePackServer struct {
	grpc.ServerStream
}

func (x *smartHTTPServiceInfoRefsReceivePackServer) Send(m *InfoRefsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _SmartHTTPService_PostUploadPack_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SmartHTTPServiceServer).PostUploadPack(&smartHTTPServicePostUploadPackServer{stream})
}

type SmartHTTPService_PostUploadPackServer interface {
	Send(*PostUploadPackResponse) error
	Recv() (*PostUploadPackRequest, error)
	grpc.ServerStream
}

type smartHTTPServicePostUploadPackServer struct {
	grpc.ServerStream
}

func (x *smartHTTPServicePostUploadPackServer) Send(m *PostUploadPackResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *smartHTTPServicePostUploadPackServer) Recv() (*PostUploadPackRequest, error) {
	m := new(PostUploadPackRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SmartHTTPService_PostReceivePack_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SmartHTTPServiceServer).PostReceivePack(&smartHTTPServicePostReceivePackServer{stream})
}

type SmartHTTPService_PostReceivePackServer interface {
	Send(*PostReceivePackResponse) error
	Recv() (*PostReceivePackRequest, error)
	grpc.ServerStream
}

type smartHTTPServicePostReceivePackServer struct {
	grpc.ServerStream
}

func (x *smartHTTPServicePostReceivePackServer) Send(m *PostReceivePackResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *smartHTTPServicePostReceivePackServer) Recv() (*PostReceivePackRequest, error) {
	m := new(PostReceivePackRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _SmartHTTPService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gitaly.SmartHTTPService",
	HandlerType: (*SmartHTTPServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "InfoRefsUploadPack",
			Handler:       _SmartHTTPService_InfoRefsUploadPack_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "InfoRefsReceivePack",
			Handler:       _SmartHTTPService_InfoRefsReceivePack_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "PostUploadPack",
			Handler:       _SmartHTTPService_PostUploadPack_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "PostReceivePack",
			Handler:       _SmartHTTPService_PostReceivePack_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "smarthttp.proto",
}

func init() { proto.RegisterFile("smarthttp.proto", fileDescriptor_smarthttp_d15d08ac1e07ff5f) }

var fileDescriptor_smarthttp_d15d08ac1e07ff5f = []byte{
	// 465 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x53, 0xcb, 0x6e, 0xd3, 0x40,
	0x14, 0x65, 0x9c, 0x34, 0xd0, 0x9b, 0x40, 0xa2, 0x5b, 0x41, 0xad, 0x48, 0xd0, 0x60, 0x24, 0xe4,
	0x45, 0xf3, 0x50, 0xd8, 0xb1, 0x84, 0x0d, 0x15, 0x48, 0x58, 0x6e, 0x23, 0x21, 0x36, 0xd6, 0xc4,
	0x9e, 0x4c, 0x47, 0x4c, 0x3c, 0xc6, 0x33, 0xad, 0xd4, 0xff, 0x40, 0x62, 0xc7, 0x77, 0xf0, 0x35,
	0x7c, 0x04, 0x5f, 0x80, 0xe2, 0x47, 0x9d, 0xc6, 0x35, 0x42, 0x20, 0x76, 0x33, 0xe7, 0x9e, 0xfb,
	0x38, 0x67, 0xee, 0x40, 0x5f, 0xaf, 0x69, 0x6a, 0xce, 0x8d, 0x49, 0x26, 0x49, 0xaa, 0x8c, 0xc2,
	0x0e, 0x17, 0x86, 0xca, 0xab, 0x61, 0x4f, 0x9f, 0xd3, 0x94, 0x45, 0x39, 0xea, 0x7c, 0x23, 0xd0,
	0x3f, 0x89, 0x57, 0xca, 0x67, 0x2b, 0xed, 0xb3, 0xcf, 0x17, 0x4c, 0x1b, 0x9c, 0x03, 0xa4, 0x2c,
	0x51, 0x5a, 0x18, 0x95, 0x5e, 0xd9, 0x64, 0x44, 0xdc, 0xee, 0x1c, 0x27, 0x79, 0xfa, 0xc4, 0xbf,
	0x8e, 0xf8, 0x5b, 0x2c, 0x3c, 0x06, 0xe4, 0xc2, 0x04, 0xa1, 0x8a, 0x57, 0x82, 0x07, 0x2a, 0x31,
	0x42, 0xc5, 0xda, 0xb6, 0x46, 0x2d, 0x77, 0xdf, 0x1f, 0x70, 0x61, 0x5e, 0x67, 0x81, 0xf7, 0x39,
	0x8e, 0x4f, 0xa1, 0xb7, 0x61, 0x67, 0x23, 0x84, 0x4a, 0xda, 0xad, 0x11, 0x71, 0xf7, 0xfd, 0x2e,
	0x17, 0xc6, 0x2b, 0xa0, 0x97, 0x9d, 0x9f, 0x5f, 0x5d, 0xeb, 0x9e, 0xe5, 0x3c, 0x87, 0x41, 0x35,
	0x9f, 0x4e, 0x54, 0xac, 0x19, 0x22, 0xb4, 0x23, 0x6a, 0x68, 0x36, 0x5a, 0xcf, 0xcf, 0xce, 0xce,
	0x77, 0x02, 0x0f, 0x3d, 0xa5, 0xcd, 0x22, 0x91, 0x8a, 0x46, 0x1e, 0x0d, 0x3f, 0xfd, 0x8b, 0x9c,
	0xb2, 0x83, 0x55, 0x75, 0x68, 0x90, 0xd8, 0xfa, 0x43, 0x89, 0xed, 0x66, 0x89, 0xc7, 0xf0, 0x68,
	0x77, 0xf2, 0xdf, 0x08, 0xfd, 0x62, 0xe5, 0x74, 0x9f, 0x85, 0x4c, 0x5c, 0xb2, 0xff, 0xa1, 0xf4,
	0x00, 0xf6, 0xb8, 0x0c, 0x44, 0x54, 0xbc, 0x4b, 0x9b, 0xcb, 0x93, 0x08, 0x9f, 0xc1, 0x7d, 0x2e,
	0x83, 0xad, 0xfa, 0xb9, 0xa2, 0x1e, 0x97, 0x55, 0x65, 0x3c, 0x82, 0x2e, 0x97, 0xc1, 0x85, 0x66,
	0x69, 0x4c, 0xd7, 0xcc, 0xde, 0xcb, 0x28, 0xc0, 0xe5, 0xa2, 0x40, 0x6a, 0xb6, 0x74, 0x6a, 0xb6,
	0x34, 0xf8, 0x7c, 0xf7, 0x76, 0x9f, 0x0b, 0x13, 0x89, 0x33, 0x86, 0xc3, 0x9a, 0x2b, 0xcd, 0x2e,
	0xce, 0x7f, 0x58, 0x30, 0x38, 0xdd, 0xfc, 0x90, 0x37, 0x67, 0x67, 0xde, 0x29, 0x4b, 0x2f, 0x45,
	0xc8, 0xf0, 0x2d, 0x60, 0xb9, 0x6b, 0xd5, 0x63, 0xe0, 0x61, 0xe9, 0xe0, 0xce, 0x3f, 0x19, 0xda,
	0xf5, 0x40, 0xde, 0xd1, 0xb9, 0x33, 0x23, 0xf8, 0x0e, 0x0e, 0x2a, 0xfc, 0x7a, 0xa8, 0xbf, 0xad,
	0xb6, 0x80, 0x07, 0x37, 0x77, 0x04, 0x1f, 0x97, 0xfc, 0x5b, 0xb7, 0x7e, 0xf8, 0xa4, 0x29, 0x5c,
	0x16, 0x75, 0xc9, 0x8c, 0xe0, 0x07, 0xe8, 0xef, 0xb8, 0x86, 0x37, 0x12, 0xeb, 0x4b, 0x36, 0x3c,
	0x6a, 0x8c, 0x6f, 0x57, 0x7e, 0x35, 0xfb, 0xb8, 0xe1, 0x49, 0xba, 0x9c, 0x84, 0x6a, 0x3d, 0xcd,
	0x8f, 0x63, 0x95, 0xf2, 0x69, 0x9e, 0x3d, 0xce, 0x36, 0x60, 0xca, 0x55, 0x71, 0x4f, 0x96, 0xcb,
	0x4e, 0x06, 0xbd, 0xf8, 0x15, 0x00, 0x00, 0xff, 0xff, 0x06, 0xb7, 0x75, 0xf3, 0xba, 0x04, 0x00,
	0x00,
}
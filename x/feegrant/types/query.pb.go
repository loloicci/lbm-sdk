// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lbm/feegrant/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	query "github.com/line/lbm-sdk/types/query"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryAllowanceRequest is the request type for the Query/Allowance RPC method.
type QueryAllowanceRequest struct {
	// granter is the address of the user granting an allowance of their funds.
	Granter string `protobuf:"bytes,1,opt,name=granter,proto3" json:"granter,omitempty" yaml:"granter_address"`
	// grantee is the address of the user being granted an allowance of another user's funds.
	Grantee string `protobuf:"bytes,2,opt,name=grantee,proto3" json:"grantee,omitempty" yaml:"grantee_address"`
}

func (m *QueryAllowanceRequest) Reset()         { *m = QueryAllowanceRequest{} }
func (m *QueryAllowanceRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllowanceRequest) ProtoMessage()    {}
func (*QueryAllowanceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2fd9dc33a9e9ee8, []int{0}
}
func (m *QueryAllowanceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllowanceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllowanceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllowanceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllowanceRequest.Merge(m, src)
}
func (m *QueryAllowanceRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllowanceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllowanceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllowanceRequest proto.InternalMessageInfo

func (m *QueryAllowanceRequest) GetGranter() string {
	if m != nil {
		return m.Granter
	}
	return ""
}

func (m *QueryAllowanceRequest) GetGrantee() string {
	if m != nil {
		return m.Grantee
	}
	return ""
}

// QueryAllowanceResponse is the response type for the Query/Allowance RPC method.
type QueryAllowanceResponse struct {
	// allowance is a allowance granted for grantee by granter.
	Allowance *Grant `protobuf:"bytes,1,opt,name=allowance,proto3" json:"allowance,omitempty"`
}

func (m *QueryAllowanceResponse) Reset()         { *m = QueryAllowanceResponse{} }
func (m *QueryAllowanceResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAllowanceResponse) ProtoMessage()    {}
func (*QueryAllowanceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2fd9dc33a9e9ee8, []int{1}
}
func (m *QueryAllowanceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllowanceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllowanceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllowanceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllowanceResponse.Merge(m, src)
}
func (m *QueryAllowanceResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllowanceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllowanceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllowanceResponse proto.InternalMessageInfo

func (m *QueryAllowanceResponse) GetAllowance() *Grant {
	if m != nil {
		return m.Allowance
	}
	return nil
}

// QueryAllowancesRequest is the request type for the Query/Allowances RPC method.
type QueryAllowancesRequest struct {
	Grantee string `protobuf:"bytes,1,opt,name=grantee,proto3" json:"grantee,omitempty" yaml:"grantee_address"`
	// pagination defines an pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllowancesRequest) Reset()         { *m = QueryAllowancesRequest{} }
func (m *QueryAllowancesRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllowancesRequest) ProtoMessage()    {}
func (*QueryAllowancesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2fd9dc33a9e9ee8, []int{2}
}
func (m *QueryAllowancesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllowancesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllowancesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllowancesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllowancesRequest.Merge(m, src)
}
func (m *QueryAllowancesRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllowancesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllowancesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllowancesRequest proto.InternalMessageInfo

func (m *QueryAllowancesRequest) GetGrantee() string {
	if m != nil {
		return m.Grantee
	}
	return ""
}

func (m *QueryAllowancesRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryAllowancesResponse is the response type for the Query/Allowances RPC method.
type QueryAllowancesResponse struct {
	// allowances are allowance's granted for grantee by granter.
	Allowances []*Grant `protobuf:"bytes,1,rep,name=allowances,proto3" json:"allowances,omitempty"`
	// pagination defines an pagination for the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllowancesResponse) Reset()         { *m = QueryAllowancesResponse{} }
func (m *QueryAllowancesResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAllowancesResponse) ProtoMessage()    {}
func (*QueryAllowancesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_e2fd9dc33a9e9ee8, []int{3}
}
func (m *QueryAllowancesResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllowancesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllowancesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllowancesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllowancesResponse.Merge(m, src)
}
func (m *QueryAllowancesResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllowancesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllowancesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllowancesResponse proto.InternalMessageInfo

func (m *QueryAllowancesResponse) GetAllowances() []*Grant {
	if m != nil {
		return m.Allowances
	}
	return nil
}

func (m *QueryAllowancesResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryAllowanceRequest)(nil), "lbm.feegrant.v1.QueryAllowanceRequest")
	proto.RegisterType((*QueryAllowanceResponse)(nil), "lbm.feegrant.v1.QueryAllowanceResponse")
	proto.RegisterType((*QueryAllowancesRequest)(nil), "lbm.feegrant.v1.QueryAllowancesRequest")
	proto.RegisterType((*QueryAllowancesResponse)(nil), "lbm.feegrant.v1.QueryAllowancesResponse")
}

func init() { proto.RegisterFile("lbm/feegrant/v1/query.proto", fileDescriptor_e2fd9dc33a9e9ee8) }

var fileDescriptor_e2fd9dc33a9e9ee8 = []byte{
	// 464 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0x31, 0x8f, 0xd3, 0x30,
	0x14, 0xc7, 0xeb, 0x22, 0x40, 0xf5, 0x0d, 0x48, 0x16, 0x94, 0x2a, 0x20, 0xdf, 0xc9, 0x12, 0x5c,
	0x97, 0xb3, 0x69, 0x39, 0xdd, 0xc0, 0x00, 0xa2, 0x0b, 0x1b, 0x82, 0x8e, 0x2c, 0xc8, 0xb9, 0x3e,
	0x4c, 0x44, 0x6a, 0xe7, 0xe2, 0xb4, 0x50, 0xa1, 0x5b, 0x60, 0x3f, 0x21, 0xe0, 0x7b, 0xf0, 0x35,
	0x18, 0x2b, 0xb1, 0x30, 0x21, 0xd4, 0xf2, 0x09, 0xf8, 0x04, 0x28, 0x4e, 0xd2, 0x54, 0x69, 0x4a,
	0x37, 0x2b, 0xff, 0xff, 0x7b, 0xef, 0xf7, 0x7f, 0x76, 0xf0, 0xad, 0xd0, 0x1f, 0x8b, 0x57, 0x00,
	0x2a, 0x96, 0x3a, 0x11, 0xd3, 0x9e, 0x38, 0x9b, 0x40, 0x3c, 0xe3, 0x51, 0x6c, 0x12, 0x43, 0xae,
	0x85, 0xfe, 0x98, 0x17, 0x22, 0x9f, 0xf6, 0xbc, 0xeb, 0xca, 0x28, 0xe3, 0x34, 0x91, 0x9e, 0x32,
	0x9b, 0x47, 0xab, 0x3d, 0x56, 0x25, 0x99, 0xce, 0x52, 0xdd, 0x97, 0x16, 0xb2, 0xe6, 0xa9, 0x23,
	0x92, 0x2a, 0xd0, 0x32, 0x09, 0x8c, 0xce, 0x3d, 0xb7, 0x95, 0x31, 0x2a, 0x04, 0x21, 0xa3, 0x40,
	0x48, 0xad, 0x4d, 0xe2, 0x44, 0x9b, 0xa9, 0xec, 0x23, 0xc2, 0x37, 0x9e, 0xa7, 0xb5, 0x8f, 0xc3,
	0xd0, 0xbc, 0x95, 0xfa, 0x14, 0x86, 0x70, 0x36, 0x01, 0x9b, 0x90, 0x63, 0x7c, 0xd5, 0x8d, 0x82,
	0xb8, 0x83, 0x0e, 0x50, 0xb7, 0x35, 0xf0, 0xfe, 0xfe, 0xda, 0x6f, 0xcf, 0xe4, 0x38, 0x7c, 0xc0,
	0x72, 0xe1, 0xa5, 0x1c, 0x8d, 0x62, 0xb0, 0x96, 0x0d, 0x0b, 0x6b, 0x59, 0x05, 0x9d, 0x66, 0x7d,
	0x15, 0x6c, 0x54, 0x01, 0x7b, 0x8a, 0xdb, 0x55, 0x08, 0x1b, 0x19, 0x6d, 0x81, 0x1c, 0xe3, 0x96,
	0x2c, 0x3e, 0x3a, 0x8e, 0xbd, 0x7e, 0x9b, 0x57, 0x96, 0xc7, 0x9f, 0xa4, 0x87, 0x61, 0x69, 0x64,
	0x17, 0xa8, 0xda, 0xd0, 0x6e, 0xc4, 0x82, 0x6d, 0xb1, 0x6a, 0x00, 0xc9, 0x43, 0x8c, 0xcb, 0xc5,
	0xba, 0x64, 0x7b, 0x7d, 0xea, 0x38, 0xd2, 0xed, 0xf3, 0xec, 0x6a, 0xa7, 0x3d, 0xfe, 0x4c, 0xaa,
	0x62, 0x81, 0xc3, 0xb5, 0x0a, 0xf6, 0x19, 0xe1, 0x9b, 0x1b, 0x40, 0x79, 0xc4, 0x13, 0x8c, 0x57,
	0xe4, 0xb6, 0x83, 0x0e, 0x2e, 0xfd, 0x27, 0xe3, 0x9a, 0x93, 0x3c, 0xaa, 0x61, 0xda, 0xdf, 0xca,
	0x94, 0x0d, 0x5b, 0x87, 0xea, 0x7f, 0x6b, 0xe2, 0xcb, 0x0e, 0x8a, 0x7c, 0x45, 0xb8, 0xb5, 0x22,
	0x23, 0x77, 0x37, 0x86, 0xd7, 0xbe, 0x10, 0xef, 0x70, 0xa7, 0x2f, 0x1b, 0xca, 0x4e, 0x3e, 0xfc,
	0xf8, 0xf3, 0xa5, 0x79, 0x8f, 0x70, 0x51, 0x7d, 0xcf, 0xab, 0x38, 0xe2, 0x7d, 0xfe, 0x82, 0xce,
	0x8b, 0x13, 0x9c, 0x93, 0x0b, 0x84, 0x71, 0xb9, 0x30, 0xb2, 0x6b, 0x5e, 0x71, 0xc7, 0x5e, 0x77,
	0xb7, 0x31, 0x27, 0x3b, 0x72, 0x64, 0x87, 0xe4, 0xce, 0x76, 0x32, 0x5b, 0x02, 0x0d, 0x06, 0xdf,
	0x17, 0x14, 0xcd, 0x17, 0x14, 0xfd, 0x5e, 0x50, 0xf4, 0x69, 0x49, 0x1b, 0xf3, 0x25, 0x6d, 0xfc,
	0x5c, 0xd2, 0xc6, 0x8b, 0xae, 0x0a, 0x92, 0xd7, 0x13, 0x9f, 0x9f, 0x9a, 0xb1, 0x08, 0x03, 0x0d,
	0x69, 0xbf, 0x23, 0x3b, 0x7a, 0x23, 0xde, 0x95, 0x5d, 0x93, 0x59, 0x04, 0xd6, 0xbf, 0xe2, 0x7e,
	0xbc, 0xfb, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x83, 0xe4, 0xf8, 0x5c, 0x20, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Allowance returns fee granted to the grantee by the granter.
	Allowance(ctx context.Context, in *QueryAllowanceRequest, opts ...grpc.CallOption) (*QueryAllowanceResponse, error)
	// Allowances returns all the grants for address.
	Allowances(ctx context.Context, in *QueryAllowancesRequest, opts ...grpc.CallOption) (*QueryAllowancesResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Allowance(ctx context.Context, in *QueryAllowanceRequest, opts ...grpc.CallOption) (*QueryAllowanceResponse, error) {
	out := new(QueryAllowanceResponse)
	err := c.cc.Invoke(ctx, "/lbm.feegrant.v1.Query/Allowance", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Allowances(ctx context.Context, in *QueryAllowancesRequest, opts ...grpc.CallOption) (*QueryAllowancesResponse, error) {
	out := new(QueryAllowancesResponse)
	err := c.cc.Invoke(ctx, "/lbm.feegrant.v1.Query/Allowances", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Allowance returns fee granted to the grantee by the granter.
	Allowance(context.Context, *QueryAllowanceRequest) (*QueryAllowanceResponse, error)
	// Allowances returns all the grants for address.
	Allowances(context.Context, *QueryAllowancesRequest) (*QueryAllowancesResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Allowance(ctx context.Context, req *QueryAllowanceRequest) (*QueryAllowanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Allowance not implemented")
}
func (*UnimplementedQueryServer) Allowances(ctx context.Context, req *QueryAllowancesRequest) (*QueryAllowancesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Allowances not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Allowance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllowanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Allowance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbm.feegrant.v1.Query/Allowance",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Allowance(ctx, req.(*QueryAllowanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Allowances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllowancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Allowances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lbm.feegrant.v1.Query/Allowances",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Allowances(ctx, req.(*QueryAllowancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "lbm.feegrant.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Allowance",
			Handler:    _Query_Allowance_Handler,
		},
		{
			MethodName: "Allowances",
			Handler:    _Query_Allowances_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lbm/feegrant/v1/query.proto",
}

func (m *QueryAllowanceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllowanceRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllowanceRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Grantee) > 0 {
		i -= len(m.Grantee)
		copy(dAtA[i:], m.Grantee)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Grantee)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Granter) > 0 {
		i -= len(m.Granter)
		copy(dAtA[i:], m.Granter)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Granter)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAllowanceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllowanceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllowanceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Allowance != nil {
		{
			size, err := m.Allowance.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAllowancesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllowancesRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllowancesRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Grantee) > 0 {
		i -= len(m.Grantee)
		copy(dAtA[i:], m.Grantee)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Grantee)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAllowancesResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllowancesResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllowancesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Allowances) > 0 {
		for iNdEx := len(m.Allowances) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Allowances[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryAllowanceRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Granter)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.Grantee)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAllowanceResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Allowance != nil {
		l = m.Allowance.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAllowancesRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Grantee)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAllowancesResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Allowances) > 0 {
		for _, e := range m.Allowances {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryAllowanceRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryAllowanceRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllowanceRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Granter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Granter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Grantee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Grantee = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryAllowanceResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryAllowanceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllowanceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Allowance", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Allowance == nil {
				m.Allowance = &Grant{}
			}
			if err := m.Allowance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryAllowancesRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryAllowancesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllowancesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Grantee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Grantee = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryAllowancesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryAllowancesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllowancesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Allowances", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Allowances = append(m.Allowances, &Grant{})
			if err := m.Allowances[len(m.Allowances)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)

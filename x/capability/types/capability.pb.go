// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: lbm/capability/v1/capability.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// Capability defines an implementation of an object capability. The index
// provided to a Capability must be globally unique.
type Capability struct {
	Index uint64 `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty" yaml:"index"`
}

func (m *Capability) Reset()      { *m = Capability{} }
func (*Capability) ProtoMessage() {}
func (*Capability) Descriptor() ([]byte, []int) {
	return fileDescriptor_e47acc4344fad2c5, []int{0}
}
func (m *Capability) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Capability) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Capability.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Capability) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Capability.Merge(m, src)
}
func (m *Capability) XXX_Size() int {
	return m.Size()
}
func (m *Capability) XXX_DiscardUnknown() {
	xxx_messageInfo_Capability.DiscardUnknown(m)
}

var xxx_messageInfo_Capability proto.InternalMessageInfo

func (m *Capability) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

// Owner defines a single capability owner. An owner is defined by the name of
// capability and the module name.
type Owner struct {
	Module string `protobuf:"bytes,1,opt,name=module,proto3" json:"module,omitempty" yaml:"module"`
	Name   string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" yaml:"name"`
}

func (m *Owner) Reset()      { *m = Owner{} }
func (*Owner) ProtoMessage() {}
func (*Owner) Descriptor() ([]byte, []int) {
	return fileDescriptor_e47acc4344fad2c5, []int{1}
}
func (m *Owner) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Owner) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Owner.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Owner) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Owner.Merge(m, src)
}
func (m *Owner) XXX_Size() int {
	return m.Size()
}
func (m *Owner) XXX_DiscardUnknown() {
	xxx_messageInfo_Owner.DiscardUnknown(m)
}

var xxx_messageInfo_Owner proto.InternalMessageInfo

// CapabilityOwners defines a set of owners of a single Capability. The set of
// owners must be unique.
type CapabilityOwners struct {
	Owners []Owner `protobuf:"bytes,1,rep,name=owners,proto3" json:"owners"`
}

func (m *CapabilityOwners) Reset()         { *m = CapabilityOwners{} }
func (m *CapabilityOwners) String() string { return proto.CompactTextString(m) }
func (*CapabilityOwners) ProtoMessage()    {}
func (*CapabilityOwners) Descriptor() ([]byte, []int) {
	return fileDescriptor_e47acc4344fad2c5, []int{2}
}
func (m *CapabilityOwners) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CapabilityOwners) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CapabilityOwners.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CapabilityOwners) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CapabilityOwners.Merge(m, src)
}
func (m *CapabilityOwners) XXX_Size() int {
	return m.Size()
}
func (m *CapabilityOwners) XXX_DiscardUnknown() {
	xxx_messageInfo_CapabilityOwners.DiscardUnknown(m)
}

var xxx_messageInfo_CapabilityOwners proto.InternalMessageInfo

func (m *CapabilityOwners) GetOwners() []Owner {
	if m != nil {
		return m.Owners
	}
	return nil
}

func init() {
	proto.RegisterType((*Capability)(nil), "lbm.capability.v1.Capability")
	proto.RegisterType((*Owner)(nil), "lbm.capability.v1.Owner")
	proto.RegisterType((*CapabilityOwners)(nil), "lbm.capability.v1.CapabilityOwners")
}

func init() {
	proto.RegisterFile("lbm/capability/v1/capability.proto", fileDescriptor_e47acc4344fad2c5)
}

var fileDescriptor_e47acc4344fad2c5 = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0xca, 0x49, 0xca, 0xd5,
	0x4f, 0x4e, 0x2c, 0x48, 0x4c, 0xca, 0xcc, 0xc9, 0x2c, 0xa9, 0xd4, 0x2f, 0x33, 0x44, 0xe2, 0xe9,
	0x15, 0x14, 0xe5, 0x97, 0xe4, 0x0b, 0x09, 0xe6, 0x24, 0xe5, 0xea, 0x21, 0x89, 0x96, 0x19, 0x4a,
	0x89, 0xa4, 0xe7, 0xa7, 0xe7, 0x83, 0x65, 0xf5, 0x41, 0x2c, 0x88, 0x42, 0x25, 0x2b, 0x2e, 0x2e,
	0x67, 0xb8, 0x32, 0x21, 0x35, 0x2e, 0xd6, 0xcc, 0xbc, 0x94, 0xd4, 0x0a, 0x09, 0x46, 0x05, 0x46,
	0x0d, 0x16, 0x27, 0x81, 0x4f, 0xf7, 0xe4, 0x79, 0x2a, 0x13, 0x73, 0x73, 0xac, 0x94, 0xc0, 0xc2,
	0x4a, 0x41, 0x10, 0x69, 0x2b, 0x96, 0x19, 0x0b, 0xe4, 0x19, 0x94, 0x12, 0xb9, 0x58, 0xfd, 0xcb,
	0xf3, 0x52, 0x8b, 0x84, 0x34, 0xb9, 0xd8, 0x72, 0xf3, 0x53, 0x4a, 0x73, 0x52, 0xc1, 0xfa, 0x38,
	0x9d, 0x04, 0x3f, 0xdd, 0x93, 0xe7, 0x85, 0xe8, 0x83, 0x88, 0x2b, 0x05, 0x41, 0x15, 0x08, 0x29,
	0x73, 0xb1, 0xe4, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x81, 0x15, 0xf2, 0x7f, 0xba, 0x27, 0xcf, 0x0d,
	0x51, 0x08, 0x12, 0x55, 0x0a, 0x02, 0x4b, 0x5a, 0x71, 0x74, 0x2c, 0x90, 0x67, 0x00, 0x5b, 0xe1,
	0xc5, 0x25, 0x80, 0x70, 0x1e, 0xd8, 0xb2, 0x62, 0x21, 0x33, 0x2e, 0xb6, 0x7c, 0x30, 0x4b, 0x82,
	0x51, 0x81, 0x59, 0x83, 0xdb, 0x48, 0x42, 0x0f, 0xc3, 0xb3, 0x7a, 0x60, 0xa5, 0x4e, 0x2c, 0x27,
	0xee, 0xc9, 0x33, 0x04, 0x41, 0x55, 0x3b, 0xb9, 0x9c, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c,
	0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1,
	0x1c, 0x43, 0x94, 0x56, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x7e, 0x4e,
	0x66, 0x5e, 0xaa, 0x7e, 0x4e, 0x52, 0xae, 0x6e, 0x71, 0x4a, 0xb6, 0x7e, 0x05, 0x72, 0x38, 0x97,
	0x54, 0x16, 0xa4, 0x16, 0x27, 0xb1, 0x81, 0xc3, 0xcd, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xce,
	0x80, 0xf8, 0x6c, 0x86, 0x01, 0x00, 0x00,
}

func (m *Capability) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Capability) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Capability) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Index != 0 {
		i = encodeVarintCapability(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Owner) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Owner) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Owner) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintCapability(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Module) > 0 {
		i -= len(m.Module)
		copy(dAtA[i:], m.Module)
		i = encodeVarintCapability(dAtA, i, uint64(len(m.Module)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *CapabilityOwners) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CapabilityOwners) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CapabilityOwners) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Owners) > 0 {
		for iNdEx := len(m.Owners) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Owners[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCapability(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintCapability(dAtA []byte, offset int, v uint64) int {
	offset -= sovCapability(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Capability) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Index != 0 {
		n += 1 + sovCapability(uint64(m.Index))
	}
	return n
}

func (m *Owner) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Module)
	if l > 0 {
		n += 1 + l + sovCapability(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovCapability(uint64(l))
	}
	return n
}

func (m *CapabilityOwners) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Owners) > 0 {
		for _, e := range m.Owners {
			l = e.Size()
			n += 1 + l + sovCapability(uint64(l))
		}
	}
	return n
}

func sovCapability(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCapability(x uint64) (n int) {
	return sovCapability(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Capability) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCapability
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
			return fmt.Errorf("proto: Capability: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Capability: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapability
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCapability(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCapability
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
func (m *Owner) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCapability
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
			return fmt.Errorf("proto: Owner: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Owner: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Module", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapability
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
				return ErrInvalidLengthCapability
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCapability
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Module = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapability
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
				return ErrInvalidLengthCapability
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCapability
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCapability(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCapability
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
func (m *CapabilityOwners) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCapability
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
			return fmt.Errorf("proto: CapabilityOwners: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CapabilityOwners: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owners", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCapability
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
				return ErrInvalidLengthCapability
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCapability
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owners = append(m.Owners, Owner{})
			if err := m.Owners[len(m.Owners)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCapability(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCapability
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
func skipCapability(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCapability
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
					return 0, ErrIntOverflowCapability
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
					return 0, ErrIntOverflowCapability
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
				return 0, ErrInvalidLengthCapability
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCapability
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCapability
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCapability        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCapability          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCapability = fmt.Errorf("proto: unexpected end of group")
)

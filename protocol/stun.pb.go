// Code generated by protoc-gen-go. DO NOT EDIT.
// source: stun.proto

package protocol

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Stun_Type int32

const (
	Stun_CONNECT                   Stun_Type = 0
	Stun_HOLE_PUNCH_REQUEST        Stun_Type = 1
	Stun_KEEP_ALIVE                Stun_Type = 2
	Stun_HOLE_PUNCH_REQUEST_ACCEPT Stun_Type = 3
)

var Stun_Type_name = map[int32]string{
	0: "CONNECT",
	1: "HOLE_PUNCH_REQUEST",
	2: "KEEP_ALIVE",
	3: "HOLE_PUNCH_REQUEST_ACCEPT",
}

var Stun_Type_value = map[string]int32{
	"CONNECT":                   0,
	"HOLE_PUNCH_REQUEST":        1,
	"KEEP_ALIVE":                2,
	"HOLE_PUNCH_REQUEST_ACCEPT": 3,
}

func (x Stun_Type) String() string {
	return proto.EnumName(Stun_Type_name, int32(x))
}

func (Stun_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d956532a5442e56c, []int{0, 0}
}

type Stun struct {
	Type                    Stun_Type                     `protobuf:"varint,1,opt,name=type,proto3,enum=protocol.Stun_Type" json:"type,omitempty"`
	HolePunchRequestMessage *Stun_HolePunchRequestMessage `protobuf:"bytes,2,opt,name=holePunchRequestMessage,proto3" json:"holePunchRequestMessage,omitempty"`
	XXX_NoUnkeyedLiteral    struct{}                      `json:"-"`
	XXX_unrecognized        []byte                        `json:"-"`
	XXX_sizecache           int32                         `json:"-"`
}

func (m *Stun) Reset()         { *m = Stun{} }
func (m *Stun) String() string { return proto.CompactTextString(m) }
func (*Stun) ProtoMessage()    {}
func (*Stun) Descriptor() ([]byte, []int) {
	return fileDescriptor_d956532a5442e56c, []int{0}
}

func (m *Stun) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stun.Unmarshal(m, b)
}
func (m *Stun) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stun.Marshal(b, m, deterministic)
}
func (m *Stun) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stun.Merge(m, src)
}
func (m *Stun) XXX_Size() int {
	return xxx_messageInfo_Stun.Size(m)
}
func (m *Stun) XXX_DiscardUnknown() {
	xxx_messageInfo_Stun.DiscardUnknown(m)
}

var xxx_messageInfo_Stun proto.InternalMessageInfo

func (m *Stun) GetType() Stun_Type {
	if m != nil {
		return m.Type
	}
	return Stun_CONNECT
}

func (m *Stun) GetHolePunchRequestMessage() *Stun_HolePunchRequestMessage {
	if m != nil {
		return m.HolePunchRequestMessage
	}
	return nil
}

type Stun_HolePunchRequestMessage struct {
	ConnectToPeerID      []byte   `protobuf:"bytes,1,opt,name=connectToPeerID,proto3" json:"connectToPeerID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stun_HolePunchRequestMessage) Reset()         { *m = Stun_HolePunchRequestMessage{} }
func (m *Stun_HolePunchRequestMessage) String() string { return proto.CompactTextString(m) }
func (*Stun_HolePunchRequestMessage) ProtoMessage()    {}
func (*Stun_HolePunchRequestMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_d956532a5442e56c, []int{0, 0}
}

func (m *Stun_HolePunchRequestMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stun_HolePunchRequestMessage.Unmarshal(m, b)
}
func (m *Stun_HolePunchRequestMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stun_HolePunchRequestMessage.Marshal(b, m, deterministic)
}
func (m *Stun_HolePunchRequestMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stun_HolePunchRequestMessage.Merge(m, src)
}
func (m *Stun_HolePunchRequestMessage) XXX_Size() int {
	return xxx_messageInfo_Stun_HolePunchRequestMessage.Size(m)
}
func (m *Stun_HolePunchRequestMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_Stun_HolePunchRequestMessage.DiscardUnknown(m)
}

var xxx_messageInfo_Stun_HolePunchRequestMessage proto.InternalMessageInfo

func (m *Stun_HolePunchRequestMessage) GetConnectToPeerID() []byte {
	if m != nil {
		return m.ConnectToPeerID
	}
	return nil
}

func init() {
	proto.RegisterEnum("protocol.Stun_Type", Stun_Type_name, Stun_Type_value)
	proto.RegisterType((*Stun)(nil), "protocol.Stun")
	proto.RegisterType((*Stun_HolePunchRequestMessage)(nil), "protocol.Stun.HolePunchRequestMessage")
}

func init() { proto.RegisterFile("stun.proto", fileDescriptor_d956532a5442e56c) }

var fileDescriptor_d956532a5442e56c = []byte{
	// 237 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x2e, 0x29, 0xcd,
	0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53, 0xc9, 0xf9, 0x39, 0x4a, 0x9b, 0x98,
	0xb8, 0x58, 0x82, 0x4b, 0x4a, 0xf3, 0x84, 0xd4, 0xb9, 0x58, 0x4a, 0x2a, 0x0b, 0x52, 0x25, 0x18,
	0x15, 0x18, 0x35, 0xf8, 0x8c, 0x84, 0xf5, 0x60, 0x2a, 0xf4, 0x40, 0xb2, 0x7a, 0x21, 0x95, 0x05,
	0xa9, 0x41, 0x60, 0x05, 0x42, 0x09, 0x5c, 0xe2, 0x19, 0xf9, 0x39, 0xa9, 0x01, 0xa5, 0x79, 0xc9,
	0x19, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0xbe, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x12,
	0x4c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0x6a, 0x68, 0x7a, 0x3d, 0xb0, 0xab, 0x0e, 0xc2, 0x65, 0x8c,
	0x94, 0x33, 0x97, 0x38, 0x0e, 0x3d, 0x42, 0x1a, 0x5c, 0xfc, 0xc9, 0xf9, 0x79, 0x79, 0xa9, 0xc9,
	0x25, 0x21, 0xf9, 0x01, 0xa9, 0xa9, 0x45, 0x9e, 0x2e, 0x60, 0x07, 0xf3, 0x04, 0xa1, 0x0b, 0x2b,
	0x45, 0x71, 0xb1, 0x80, 0x1c, 0x2d, 0xc4, 0xcd, 0xc5, 0xee, 0xec, 0xef, 0xe7, 0xe7, 0xea, 0x1c,
	0x22, 0xc0, 0x20, 0x24, 0xc6, 0x25, 0xe4, 0xe1, 0xef, 0xe3, 0x1a, 0x1f, 0x10, 0xea, 0xe7, 0xec,
	0x11, 0x1f, 0xe4, 0x1a, 0x18, 0xea, 0x1a, 0x1c, 0x22, 0xc0, 0x28, 0xc4, 0xc7, 0xc5, 0xe5, 0xed,
	0xea, 0x1a, 0x10, 0xef, 0xe8, 0xe3, 0x19, 0xe6, 0x2a, 0xc0, 0x24, 0x24, 0xcb, 0x25, 0x89, 0xa9,
	0x2e, 0xde, 0xd1, 0xd9, 0xd9, 0x35, 0x20, 0x44, 0x80, 0x39, 0x89, 0x0d, 0xec, 0x41, 0x63, 0x40,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x8a, 0x2f, 0x55, 0xac, 0x53, 0x01, 0x00, 0x00,
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test.proto

package test

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Test struct {
	//姓名
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	//体重
	Tizhong []int32 `protobuf:"varint,2,rep,packed,name=tizhong,proto3" json:"tizhong,omitempty"`
	//身高
	Shengao int32 `protobuf:"varint,3,opt,name=shengao,proto3" json:"shengao,omitempty"`
	//格言
	Motto                string   `protobuf:"bytes,4,opt,name=motto,proto3" json:"motto,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Test) Reset()         { *m = Test{} }
func (m *Test) String() string { return proto.CompactTextString(m) }
func (*Test) ProtoMessage()    {}
func (*Test) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{0}
}

func (m *Test) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Test.Unmarshal(m, b)
}
func (m *Test) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Test.Marshal(b, m, deterministic)
}
func (m *Test) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Test.Merge(m, src)
}
func (m *Test) XXX_Size() int {
	return xxx_messageInfo_Test.Size(m)
}
func (m *Test) XXX_DiscardUnknown() {
	xxx_messageInfo_Test.DiscardUnknown(m)
}

var xxx_messageInfo_Test proto.InternalMessageInfo

func (m *Test) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Test) GetTizhong() []int32 {
	if m != nil {
		return m.Tizhong
	}
	return nil
}

func (m *Test) GetShengao() int32 {
	if m != nil {
		return m.Shengao
	}
	return 0
}

func (m *Test) GetMotto() string {
	if m != nil {
		return m.Motto
	}
	return ""
}

func init() {
	proto.RegisterType((*Test)(nil), "test.Test")
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_c161fcfdc0c3ff1e) }

var fileDescriptor_c161fcfdc0c3ff1e = []byte{
	// 121 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x49, 0x2d, 0x2e,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0x52, 0xb8, 0x58, 0x42, 0x52,
	0x8b, 0x4b, 0x84, 0x84, 0xb8, 0x58, 0xf2, 0x12, 0x73, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38,
	0x83, 0xc0, 0x6c, 0x21, 0x09, 0x2e, 0xf6, 0x92, 0xcc, 0xaa, 0x8c, 0xfc, 0xbc, 0x74, 0x09, 0x26,
	0x05, 0x66, 0x0d, 0xd6, 0x20, 0x18, 0x17, 0x24, 0x53, 0x9c, 0x91, 0x9a, 0x97, 0x9e, 0x98, 0x2f,
	0xc1, 0xac, 0xc0, 0x08, 0x92, 0x81, 0x72, 0x85, 0x44, 0xb8, 0x58, 0x73, 0xf3, 0x4b, 0x4a, 0xf2,
	0x25, 0x58, 0xc0, 0x06, 0x41, 0x38, 0x49, 0x6c, 0x60, 0x2b, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xb1, 0x07, 0x65, 0xd4, 0x80, 0x00, 0x00, 0x00,
}

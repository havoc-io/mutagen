// Code generated by protoc-gen-go. DO NOT EDIT.
// source: synchronization/core/change.proto

package core

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

type Change struct {
	Path                 string   `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Old                  *Entry   `protobuf:"bytes,2,opt,name=old,proto3" json:"old,omitempty"`
	New                  *Entry   `protobuf:"bytes,3,opt,name=new,proto3" json:"new,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Change) Reset()         { *m = Change{} }
func (m *Change) String() string { return proto.CompactTextString(m) }
func (*Change) ProtoMessage()    {}
func (*Change) Descriptor() ([]byte, []int) {
	return fileDescriptor_abf1663c97210f4e, []int{0}
}

func (m *Change) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Change.Unmarshal(m, b)
}
func (m *Change) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Change.Marshal(b, m, deterministic)
}
func (m *Change) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Change.Merge(m, src)
}
func (m *Change) XXX_Size() int {
	return xxx_messageInfo_Change.Size(m)
}
func (m *Change) XXX_DiscardUnknown() {
	xxx_messageInfo_Change.DiscardUnknown(m)
}

var xxx_messageInfo_Change proto.InternalMessageInfo

func (m *Change) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Change) GetOld() *Entry {
	if m != nil {
		return m.Old
	}
	return nil
}

func (m *Change) GetNew() *Entry {
	if m != nil {
		return m.New
	}
	return nil
}

func init() {
	proto.RegisterType((*Change)(nil), "core.Change")
}

func init() { proto.RegisterFile("synchronization/core/change.proto", fileDescriptor_abf1663c97210f4e) }

var fileDescriptor_abf1663c97210f4e = []byte{
	// 174 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2c, 0xae, 0xcc, 0x4b,
	0xce, 0x28, 0xca, 0xcf, 0xcb, 0xac, 0x4a, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x4f, 0xce, 0x2f, 0x4a,
	0xd5, 0x4f, 0xce, 0x48, 0xcc, 0x4b, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01,
	0x09, 0x49, 0x29, 0x60, 0x55, 0x98, 0x9a, 0x57, 0x52, 0x54, 0x09, 0x51, 0xa7, 0x14, 0xc5, 0xc5,
	0xe6, 0x0c, 0xd6, 0x27, 0x24, 0xc4, 0xc5, 0x52, 0x90, 0x58, 0x92, 0x21, 0xc1, 0xa8, 0xc0, 0xa8,
	0xc1, 0x19, 0x04, 0x66, 0x0b, 0xc9, 0x72, 0x31, 0xe7, 0xe7, 0xa4, 0x48, 0x30, 0x29, 0x30, 0x6a,
	0x70, 0x1b, 0x71, 0xeb, 0x81, 0x74, 0xeb, 0xb9, 0x82, 0x74, 0x07, 0x81, 0xc4, 0x41, 0xd2, 0x79,
	0xa9, 0xe5, 0x12, 0xcc, 0x58, 0xa4, 0xf3, 0x52, 0xcb, 0x9d, 0x2c, 0xa2, 0xcc, 0xd2, 0x33, 0x4b,
	0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0x73, 0x4b, 0x4b, 0x12, 0xd3, 0x53, 0xf3, 0x74,
	0x33, 0xf3, 0x61, 0x4c, 0xfd, 0x82, 0xec, 0x74, 0x7d, 0x6c, 0x2e, 0x4c, 0x62, 0x03, 0x3b, 0xce,
	0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xc3, 0xfc, 0xe4, 0xdb, 0xe9, 0x00, 0x00, 0x00,
}
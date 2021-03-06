// Code generated by protoc-gen-go.
// source: buss.proto
// DO NOT EDIT!

/*
Package buss is a generated protocol buffer package.

It is generated from these files:
	buss.proto

It has these top-level messages:
	Pull
	PushReq
	PushRes
*/
package buss

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ResultCode int32

const (
	ResultCode_CLIENT_RESPONSE_OK  ResultCode = 0
	ResultCode_CLIENT_IN_OTHER_SVR ResultCode = 1
	ResultCode_CLIENT_NOT_EXSIT    ResultCode = 2
)

var ResultCode_name = map[int32]string{
	0: "CLIENT_RESPONSE_OK",
	1: "CLIENT_IN_OTHER_SVR",
	2: "CLIENT_NOT_EXSIT",
}
var ResultCode_value = map[string]int32{
	"CLIENT_RESPONSE_OK":  0,
	"CLIENT_IN_OTHER_SVR": 1,
	"CLIENT_NOT_EXSIT":    2,
}

func (x ResultCode) Enum() *ResultCode {
	p := new(ResultCode)
	*p = x
	return p
}
func (x ResultCode) String() string {
	return proto.EnumName(ResultCode_name, int32(x))
}
func (x *ResultCode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ResultCode_value, data, "ResultCode")
	if err != nil {
		return err
	}
	*x = ResultCode(value)
	return nil
}
func (ResultCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Pull struct {
	ConnecId         *int32 `protobuf:"varint,1,req,name=connec_id" json:"connec_id,omitempty"`
	IoThreadId       *int32 `protobuf:"varint,2,req,name=io_thread_id" json:"io_thread_id,omitempty"`
	LinkId           *int64 `protobuf:"varint,3,req,name=link_id" json:"link_id,omitempty"`
	UserData         []byte `protobuf:"bytes,4,opt,name=user_data" json:"user_data,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Pull) Reset()                    { *m = Pull{} }
func (m *Pull) String() string            { return proto.CompactTextString(m) }
func (*Pull) ProtoMessage()               {}
func (*Pull) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Pull) GetConnecId() int32 {
	if m != nil && m.ConnecId != nil {
		return *m.ConnecId
	}
	return 0
}

func (m *Pull) GetIoThreadId() int32 {
	if m != nil && m.IoThreadId != nil {
		return *m.IoThreadId
	}
	return 0
}

func (m *Pull) GetLinkId() int64 {
	if m != nil && m.LinkId != nil {
		return *m.LinkId
	}
	return 0
}

func (m *Pull) GetUserData() []byte {
	if m != nil {
		return m.UserData
	}
	return nil
}

type PushReq struct {
	LinkId           *int64 `protobuf:"varint,1,req,name=link_id" json:"link_id,omitempty"`
	UserData         []byte `protobuf:"bytes,2,req,name=user_data" json:"user_data,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PushReq) Reset()                    { *m = PushReq{} }
func (m *PushReq) String() string            { return proto.CompactTextString(m) }
func (*PushReq) ProtoMessage()               {}
func (*PushReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *PushReq) GetLinkId() int64 {
	if m != nil && m.LinkId != nil {
		return *m.LinkId
	}
	return 0
}

func (m *PushReq) GetUserData() []byte {
	if m != nil {
		return m.UserData
	}
	return nil
}

type PushRes struct {
	ResultCode       *ResultCode `protobuf:"varint,1,req,name=result_code,enum=buss.ResultCode" json:"result_code,omitempty"`
	LinkId           *int64      `protobuf:"varint,2,req,name=link_id" json:"link_id,omitempty"`
	UserData         []byte      `protobuf:"bytes,3,req,name=user_data" json:"user_data,omitempty"`
	BussSvrAddress   *string     `protobuf:"bytes,4,opt,name=buss_svr_address" json:"buss_svr_address,omitempty"`
	BussPort         *int32      `protobuf:"varint,5,opt,name=buss_port" json:"buss_port,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *PushRes) Reset()                    { *m = PushRes{} }
func (m *PushRes) String() string            { return proto.CompactTextString(m) }
func (*PushRes) ProtoMessage()               {}
func (*PushRes) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PushRes) GetResultCode() ResultCode {
	if m != nil && m.ResultCode != nil {
		return *m.ResultCode
	}
	return ResultCode_CLIENT_RESPONSE_OK
}

func (m *PushRes) GetLinkId() int64 {
	if m != nil && m.LinkId != nil {
		return *m.LinkId
	}
	return 0
}

func (m *PushRes) GetUserData() []byte {
	if m != nil {
		return m.UserData
	}
	return nil
}

func (m *PushRes) GetBussSvrAddress() string {
	if m != nil && m.BussSvrAddress != nil {
		return *m.BussSvrAddress
	}
	return ""
}

func (m *PushRes) GetBussPort() int32 {
	if m != nil && m.BussPort != nil {
		return *m.BussPort
	}
	return 0
}

func init() {
	proto.RegisterType((*Pull)(nil), "buss.Pull")
	proto.RegisterType((*PushReq)(nil), "buss.PushReq")
	proto.RegisterType((*PushRes)(nil), "buss.PushRes")
	proto.RegisterEnum("buss.ResultCode", ResultCode_name, ResultCode_value)
}

func init() { proto.RegisterFile("buss.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x90, 0x51, 0x4b, 0xc3, 0x30,
	0x14, 0x85, 0x6d, 0xd7, 0x32, 0x76, 0x2d, 0x1a, 0xa3, 0x68, 0x1f, 0x47, 0x41, 0x18, 0x82, 0x7b,
	0xf0, 0x2f, 0x8c, 0x82, 0x45, 0x69, 0x4b, 0x53, 0xc4, 0xb7, 0x4b, 0x6d, 0x02, 0x2b, 0x96, 0x66,
	0x26, 0xad, 0xef, 0xfe, 0x73, 0x93, 0xb0, 0x41, 0x61, 0x6f, 0xe1, 0xbb, 0xe7, 0x7c, 0x1c, 0x02,
	0xf0, 0x35, 0x69, 0xbd, 0x3d, 0x28, 0x39, 0x4a, 0x1a, 0xd8, 0x77, 0xc2, 0x20, 0x28, 0xa7, 0xbe,
	0xa7, 0x37, 0xb0, 0x6a, 0xe5, 0x30, 0x88, 0x16, 0x3b, 0x1e, 0x7b, 0x6b, 0x7f, 0x13, 0xd2, 0x3b,
	0x88, 0x3a, 0x89, 0xe3, 0x5e, 0x89, 0x86, 0x5b, 0xea, 0x3b, 0x7a, 0x0d, 0xcb, 0xbe, 0x1b, 0xbe,
	0x2d, 0x58, 0x18, 0xb0, 0xb0, 0xcd, 0x49, 0x0b, 0x85, 0xbc, 0x19, 0x9b, 0x38, 0x58, 0x7b, 0x9b,
	0x28, 0x79, 0x86, 0x65, 0x39, 0xe9, 0x7d, 0x25, 0x7e, 0xe6, 0x71, 0xef, 0x3c, 0x6e, 0x95, 0x51,
	0xf2, 0xe7, 0x9d, 0xf2, 0x9a, 0x3e, 0xc2, 0xa5, 0x12, 0x7a, 0xea, 0x47, 0x6c, 0x25, 0x17, 0xae,
	0x73, 0xf5, 0x42, 0xb6, 0x6e, 0x77, 0xe5, 0x0e, 0x3b, 0xc3, 0xe7, 0x5a, 0xff, 0x5c, 0x6b, 0x87,
	0x45, 0x34, 0x06, 0x62, 0x6b, 0xa8, 0x7f, 0x15, 0x36, 0x9c, 0x1b, 0xad, 0x76, 0xfb, 0x56, 0x36,
	0xec, 0x2e, 0x07, 0xa9, 0xc6, 0x38, 0x34, 0x28, 0x7c, 0x62, 0x00, 0x33, 0xfd, 0x3d, 0xd0, 0xdd,
	0x7b, 0x96, 0xe6, 0x35, 0x56, 0x29, 0x2b, 0x8b, 0x9c, 0xa5, 0x58, 0xbc, 0x91, 0x0b, 0xfa, 0x00,
	0xb7, 0x47, 0x9e, 0xe5, 0x58, 0xd4, 0xaf, 0x69, 0x85, 0xec, 0xa3, 0x22, 0x9e, 0xf9, 0x2b, 0x72,
	0x3c, 0xe4, 0x45, 0x8d, 0xe9, 0x27, 0xcb, 0x6a, 0xe2, 0xff, 0x07, 0x00, 0x00, 0xff, 0xff, 0x06,
	0xde, 0x8c, 0x95, 0x6f, 0x01, 0x00, 0x00,
}

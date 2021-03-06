// Code generated by protoc-gen-go.
// source: log_in.proto
// DO NOT EDIT!

package game

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ErrCode int32

const (
	ErrCode_AUTHEN_OK        ErrCode = 0
	ErrCode_USER_LOCKED      ErrCode = 1
	ErrCode_AREA_NOT_ALLOWED ErrCode = 2
)

var ErrCode_name = map[int32]string{
	0: "AUTHEN_OK",
	1: "USER_LOCKED",
	2: "AREA_NOT_ALLOWED",
}
var ErrCode_value = map[string]int32{
	"AUTHEN_OK":        0,
	"USER_LOCKED":      1,
	"AREA_NOT_ALLOWED": 2,
}

func (x ErrCode) Enum() *ErrCode {
	p := new(ErrCode)
	*p = x
	return p
}
func (x ErrCode) String() string {
	return proto.EnumName(ErrCode_name, int32(x))
}
func (x *ErrCode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(ErrCode_value, data, "ErrCode")
	if err != nil {
		return err
	}
	*x = ErrCode(value)
	return nil
}
func (ErrCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

type LoginRequest struct {
	UserName         *string   `protobuf:"bytes,1,req,name=user_name" json:"user_name,omitempty"`
	PassWord         *string   `protobuf:"bytes,2,req,name=pass_word" json:"pass_word,omitempty"`
	EntryptType      *int32    `protobuf:"varint,3,req,name=entrypt_type" json:"entrypt_type,omitempty"`
	Addr             *ConnAddr `protobuf:"bytes,4,opt,name=addr" json:"addr,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *LoginRequest) Reset()                    { *m = LoginRequest{} }
func (m *LoginRequest) String() string            { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()               {}
func (*LoginRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *LoginRequest) GetUserName() string {
	if m != nil && m.UserName != nil {
		return *m.UserName
	}
	return ""
}

func (m *LoginRequest) GetPassWord() string {
	if m != nil && m.PassWord != nil {
		return *m.PassWord
	}
	return ""
}

func (m *LoginRequest) GetEntryptType() int32 {
	if m != nil && m.EntryptType != nil {
		return *m.EntryptType
	}
	return 0
}

func (m *LoginRequest) GetAddr() *ConnAddr {
	if m != nil {
		return m.Addr
	}
	return nil
}

type LoginResponse struct {
	ErrCode          *ErrCode `protobuf:"varint,1,req,name=err_code,enum=game.ErrCode" json:"err_code,omitempty"`
	PlayerId         *int64   `protobuf:"varint,2,req,name=player_id" json:"player_id,omitempty"`
	ErrString        *string  `protobuf:"bytes,3,req,name=err_string" json:"err_string,omitempty"`
	GameServer       *string  `protobuf:"bytes,4,req,name=game_server" json:"game_server,omitempty"`
	Port             *string  `protobuf:"bytes,5,req,name=port" json:"port,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *LoginResponse) Reset()                    { *m = LoginResponse{} }
func (m *LoginResponse) String() string            { return proto.CompactTextString(m) }
func (*LoginResponse) ProtoMessage()               {}
func (*LoginResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *LoginResponse) GetErrCode() ErrCode {
	if m != nil && m.ErrCode != nil {
		return *m.ErrCode
	}
	return ErrCode_AUTHEN_OK
}

func (m *LoginResponse) GetPlayerId() int64 {
	if m != nil && m.PlayerId != nil {
		return *m.PlayerId
	}
	return 0
}

func (m *LoginResponse) GetErrString() string {
	if m != nil && m.ErrString != nil {
		return *m.ErrString
	}
	return ""
}

func (m *LoginResponse) GetGameServer() string {
	if m != nil && m.GameServer != nil {
		return *m.GameServer
	}
	return ""
}

func (m *LoginResponse) GetPort() string {
	if m != nil && m.Port != nil {
		return *m.Port
	}
	return ""
}

func init() {
	proto.RegisterType((*LoginRequest)(nil), "game.LoginRequest")
	proto.RegisterType((*LoginResponse)(nil), "game.LoginResponse")
	proto.RegisterEnum("game.ErrCode", ErrCode_name, ErrCode_value)
}

func init() { proto.RegisterFile("log_in.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 276 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x34, 0x8f, 0x41, 0x4b, 0xc3, 0x30,
	0x18, 0x86, 0x5d, 0xd7, 0xa1, 0xfd, 0xda, 0x6e, 0x35, 0x7a, 0x28, 0x22, 0x38, 0x76, 0x1a, 0x1e,
	0x7a, 0xd8, 0x1f, 0x90, 0x52, 0x0b, 0xc2, 0xca, 0x0a, 0x75, 0xc3, 0x63, 0x28, 0x6b, 0xa8, 0x85,
	0x99, 0xc4, 0x2f, 0x99, 0x52, 0x7f, 0xbd, 0x49, 0xab, 0xc7, 0x3c, 0xe4, 0x7d, 0x9f, 0xf7, 0x83,
	0xe0, 0x24, 0x5a, 0xda, 0xf1, 0x44, 0xa2, 0xd0, 0x82, 0xb8, 0x6d, 0xfd, 0xc1, 0xee, 0x16, 0x47,
	0xc1, 0x39, 0xad, 0x9b, 0x06, 0x47, 0xbc, 0x7a, 0x87, 0xa0, 0x10, 0x6d, 0xc7, 0x2b, 0xf6, 0x79,
	0x66, 0x4a, 0x93, 0x6b, 0xf0, 0xce, 0x8a, 0x21, 0xe5, 0xe6, 0x77, 0x3c, 0x59, 0x3a, 0x6b, 0xcf,
	0x22, 0x59, 0x2b, 0x45, 0xbf, 0x05, 0x36, 0xb1, 0x33, 0xa0, 0x5b, 0x08, 0x18, 0xd7, 0xd8, 0x4b,
	0x4d, 0x75, 0x2f, 0x59, 0x3c, 0x35, 0x74, 0x46, 0xee, 0xc1, 0xb5, 0xcd, 0xb1, 0xbb, 0x9c, 0xac,
	0xfd, 0xcd, 0x3c, 0xb1, 0xc6, 0x24, 0x33, 0xc2, 0xd4, 0xd0, 0xd5, 0x0f, 0x84, 0x7f, 0x26, 0x25,
	0x05, 0x57, 0x8c, 0x3c, 0xc0, 0x15, 0x43, 0xa4, 0x47, 0xd1, 0x8c, 0xa6, 0xf9, 0x26, 0x1c, 0x23,
	0x39, 0x62, 0x66, 0xe0, 0x20, 0x3e, 0xd5, 0xbd, 0x59, 0xd3, 0x8d, 0xe2, 0x29, 0x21, 0x00, 0x36,
	0xa3, 0x34, 0x76, 0xbc, 0x1d, 0xb4, 0x1e, 0xb9, 0x01, 0xdf, 0xc6, 0xa8, 0x99, 0xfd, 0xc5, 0xac,
	0xdd, 0xc2, 0x00, 0x5c, 0x29, 0x50, 0xc7, 0x33, 0xfb, 0x7a, 0x7c, 0x82, 0xcb, 0xff, 0xd2, 0x10,
	0xbc, 0xf4, 0xb0, 0x7f, 0xc9, 0x77, 0xb4, 0xdc, 0x46, 0x17, 0x64, 0x01, 0xfe, 0xe1, 0x35, 0xaf,
	0x68, 0x51, 0x66, 0xdb, 0xfc, 0x39, 0x9a, 0x98, 0xd3, 0xa2, 0xb4, 0xca, 0x53, 0xba, 0x2b, 0xf7,
	0x34, 0x2d, 0x8a, 0xf2, 0xcd, 0x50, 0xe7, 0x37, 0x00, 0x00, 0xff, 0xff, 0xb8, 0xd8, 0x72, 0x17,
	0x4c, 0x01, 0x00, 0x00,
}

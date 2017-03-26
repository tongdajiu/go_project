package game_core

import (
	"unsafe"
)

type PacketHeader struct {
	Check_sum      int32
	Protol_name    [12]byte
	Enctrpt_method int32
	Size           int32
}

func PackageRes(packet_req *PacketHeader, data []byte) []byte {
	var packet_buf []byte

	size := len(data) + int(unsafe.Sizeof(PacketHeader{}))
	packet_buf = make([]byte, size)
	packet_res := (*PacketHeader)(unsafe.Pointer(&packet_buf[0]))
	*packet_res = *packet_req
	copy(packet_buf[unsafe.Sizeof(PacketHeader{}):], data)
	packet_res.Size = int32(len(data))
	packet_res.Check_sum = CheckSum(packet_buf[unsafe.Sizeof(PacketHeader{}.Check_sum):])

	return packet_buf
}

func CheckSum(buf []byte) int32 {
	var check_sum int32

	check_sum = 0
	for _, v := range buf {
		check_sum += int32(v)
	}

	return check_sum
}

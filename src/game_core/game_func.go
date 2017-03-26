package game_core

import (
	"unsafe"
)

func GameParse(recv_buf []byte) ([]byte, int) {
	var buf []byte
	if uintptr(len(buf)) >= unsafe.Sizeof(PacketHeader{}) { // 收到完整的包头

		// 解析头部的包
		packet_header := (*PacketHeader)(unsafe.Pointer(&buf[0]))
		size := int(packet_header.Size) + int(unsafe.Sizeof(PacketHeader{}))

		// 收到完成的包
		if len(buf) >= size {
			buf = make([]byte, size)
			copy(buf, recv_buf[:size])

			return buf, size
		}
	}

	return buf, 0
}

func CreateTokenNums() string {
	token_nums := "abcd"

	return token_nums
}

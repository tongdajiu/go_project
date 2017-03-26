package main

import (
	"conn_link_test/buss/pb"
	"fmt"
	"net"
	"unsafe"

	"github.com/golang/protobuf/proto"
)

func Send(conn net.Conn, buf []byte) bool {
	var offset = 0
	for {
		if offset >= len(buf) {
			return true
		}

		size, err := conn.Write(buf[offset:])
		if err != nil {
			fmt.Printf("Send Failed!ErrString=%s", err.Error())
			return false
		}
		offset += size
	}
}

func Recv(conn net.Conn) ([]byte, bool) {

	recv_buf := make([]byte, 1024)
	packet_buf := make([]byte, 0, 1024)

	for {
		size, err := conn.Read(recv_buf)
		if err != nil {
			fmt.Printf("Read Failed!ErrString=%s", err.Error())
			return nil, false
		}

		packet_buf = append(packet_buf, recv_buf[:size]...)
		if len(packet_buf) >= int(unsafe.Sizeof(int32(0))) {
			packet_len := int(*(*int32)(unsafe.Pointer(&packet_buf[0])))

			if len(packet_buf) >= packet_len+int(unsafe.Sizeof(int32(0))) {
				packet_data := packet_buf[0 : packet_len+int(unsafe.Sizeof(int32(0)))]
				return packet_data, true
			}
		}
	}
}

func Run(conn net.Conn) {

	defer conn.Close()

	var packet_data []byte
	var ok bool
	var message_name string
	var body []byte

	// 发送buss认证信息
	packet_data = Pack(nil, "Confirm", 0)
	if !Send(conn, packet_data) {
		return
	}

	packet_data, ok = Recv(conn)
	if !ok {
		return
	}

	body, message_name, ok = Unpack(packet_data)
	if message_name != "Confirm" {
		fmt.Printf("Confirm Unsucceed!")
		return
	}

	fmt.Printf("Confirm Succeed!")

	for {
		packet_data, ok = Recv(conn)
		if !ok {
			return
		}

		body, message_name, ok = Unpack(packet_data)
		if message_name != "Pull" {
			fmt.Printf("Pull Unsucceed!")
			return
		}
		var pull_req buss.Pull
		var pull_res buss.Pull
		var user_data string

		proto.Unmarshal(body, &pull_req)

		fmt.Printf("Data Recved!MsgName=%s,UserData=%s, ConnectId=%d,IoThreadId=%d",
			message_name,
			*pull_req.UserData,
			*pull_req.ConnecId,
			*pull_req.IoThreadId)

		user_data = *pull_req.UserData
		user_data += "tongfw123456"

		pull_res.ConnecId = proto.Int32(*pull_req.ConnecId)
		pull_res.IoThreadId = proto.Int32(*pull_req.IoThreadId)
		pull_res.LinkId = proto.Int64(*pull_req.LinkId)
		pull_res.UserData = proto.String(user_data)

		body, _ = proto.Marshal(&pull_res)
		packet_data = Pack(body, "Pull", 0)

		if !Send(conn, packet_data) {
			return
		}

	}

}

func main() {
	var conn_nums = 10
	conn_slice := make([]net.Conn, 0, conn_nums)
	close_slice := make(chan int, conn_nums)

	defer func() {
		for _, v := range conn_slice {
			v.Close()
		}
	}()

	for i := 0; i < conn_nums; i++ {
		conn, err := net.Dial("tcp4", "192.168.170.14:12307")
		if err != nil {
			fmt.Printf("Connect Failed!ErrString=%s\n", err.Error())
			return
		}

		conn_slice = append(conn_slice, conn)
		go Run(conn)
	}

	var count = 0
	for {
		select {
		case <-close_slice:
			count++
			if count >= conn_nums {
				return
			}
			break
		}

	}
}

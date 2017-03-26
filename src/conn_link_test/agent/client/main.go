package main

import (
	"fmt"
	"net"
	"time"
)

/*
packet_len int16
type Header struct{
	name_len int16
	message_name [name_len]byte
	user_data [packet_len - unsafe.sizeof(int16) - name_len]byte
}
*/

func Run(conn net.Conn) {

	defer conn.Close()

	var size int

	// 发送验证包
	message_name := "AUTH"
	send_buf := Pack(nil, message_name)
	_, err := conn.Write(send_buf)
	if err != nil {
		fmt.Printf("Confim-Packet Faield!ErrString=%s\n", err.Error())
		return
	}

	recv_buf := make([]byte, 256)
	size, err = conn.Read(recv_buf)
	if err != nil {
		fmt.Printf("Confim-Packet Faield!ErrString=%s\n", err.Error())
		return
	}

	_, message_name, _ = Unpack(recv_buf[0:size])
	if message_name != "AUTH" {
		return
	}

	for {
		size := 0
		message_name = "ECHO"
		//abc:= fmt.Sprintf("hello world%d", time.Now().Unix())
		abc := "hello world"
		send_buf = Pack([]byte(abc), message_name)

		_, err = conn.Write(send_buf)
		if err != nil {
			fmt.Printf("Write Failed!ErrString=%s\n", err.Error())
			return
		}

		//recv_buf = make([]byte, 256)
		size, err := conn.Read(recv_buf)
		if err != nil {
			fmt.Printf("Read Failed!Errstring=%s\n", err.Error())
			return
		}

		var packet_body []byte
		var ok bool

		packet_body, message_name, ok = Unpack(recv_buf[:size])
		if !ok {
			fmt.Printf("Packet Err!")
			return
		}

		fmt.Print("************************************\n")
		fmt.Printf("Echo!String=%s\n", string(packet_body))
		fmt.Print("************************************\n")

		time.Sleep(2 * time.Second)

	}
}

func main() {
	conn, err := net.Dial("tcp4", "192.168.170.14:12306")
	if err != nil {
		fmt.Printf("Connect Failed!ErrString=%s", err.Error())
		return
	}

	Run(conn)
}

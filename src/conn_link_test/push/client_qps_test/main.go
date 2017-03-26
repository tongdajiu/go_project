package main

import (
	"conn_link_test/push/comm"
	"fmt"
	"logger"
	"math/rand"
	"net"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

/*
packet_len int16
type Header struct{
	name_len int16
	message_name [name_len]byte
	user_data [packet_len - unsafe.sizeof(int16) - name_len]byte
}
*/

func Run(stop_chann chan int) {

	defer func() {
		stop_chann <- 1
	}()

	log_obj := logger.Instance()
	var size int
	var auth_req comm.AuthReq
	var slice reflect.SliceHeader

	// 连接服务器
	remote_ip := os.Args[1]
	remote_ip += ":"
	remote_ip += os.Args[2]
	conn, err := net.DialTimeout("tcp", remote_ip, 10*time.Second)
	if err != nil {
		log_obj.LogAppError("Connect Failed!ErrString=%s", err.Error())
		return
	}

	defer conn.Close()

	// 发送验证包
	copy(auth_req.Password[:], []byte("123232"))
	copy(auth_req.User_name[:], []byte("tongfangwei"))
	auth_req.Uuid = rand.Int63n(5000) + 1

	slice.Cap = int(unsafe.Sizeof(auth_req))
	slice.Len = int(unsafe.Sizeof(auth_req))
	slice.Data = uintptr(unsafe.Pointer(&auth_req))

	buf := *(*[]byte)(unsafe.Pointer(&slice))

	message_name := "AUTH"
	send_buf := comm.Pack(buf, message_name)
	_, err = conn.Write(send_buf)
	if err != nil {
		log_obj.LogAppError("Confim-Packet Faield!ErrString=%s\n", err.Error())
		return
	}

	// 设置超时
	recv_buf := make([]byte, 256)
	//conn.SetReadDeadline(time.Now().Add(10 * time.Second)) // 超时10s
	size, err = conn.Read(recv_buf)
	if err != nil {
		log_obj.LogAppError("Confim-Packet Faield!ErrString=%s\n", err.Error())
		return
	}

	_, message_name, _ = comm.Unpack(recv_buf[0:size])
	if message_name != "AUTH" {
		return
	}

	// 接受推送数据
	for {
		var push_req *comm.UserPushReq
		var push_res comm.UserPushRes

		size, err = conn.Read(recv_buf[:])
		if err != nil {
			log_obj.LogAppError("Push Data Recv Failed!ErrString=%s\n", err.Error())
			return
		}
		buf, message_name, _ = comm.Unpack(recv_buf[:size])
		push_req = (*comm.UserPushReq)(unsafe.Pointer(&buf[0]))

		// 应答
		push_res.Content_id = push_req.Content_id
		push_res.Stat = 0

		slice.Cap = int(unsafe.Sizeof(push_res))
		slice.Len = int(unsafe.Sizeof(push_res))
		slice.Data = uintptr(unsafe.Pointer(&push_res))

		buf = *(*[]byte)(unsafe.Pointer(&slice))

		send_buf = comm.Pack(buf, "NOTICE")

		_, err = conn.Write(send_buf)
		if err != nil {
			log_obj.LogAppError("Push Data Send Failed!ErrString=%s\n", err.Error())
			return
		}
	}
}

func main() {
	log_obj := logger.Instance()
	if err := log_obj.Load("./conf/conf.xml"); err != nil {
		fmt.Printf("Load Xml Failed!ErrString=%s\n", err.Error())
		return
	}

	defer log_obj.Close()

	arg_nums := len(os.Args)
	if arg_nums < 4 {
		log_obj.LogAppError("ArgNums Not Right!ArgNums=%d", len(os.Args))
		return
	}

	var conn_nums int

	rand.Seed(time.Now().Unix())
	conn_nums, _ = strconv.Atoi(os.Args[3])
	stop_chann := make(chan int, 1024)
	for i := 0; i < conn_nums; i++ {
		go Run(stop_chann)
	}

	nums := 0
	for {
		select {
		case <-stop_chann:
			nums++
			break
		}

		if nums >= conn_nums {
			log_obj.LogAppInfo("Test Over")
			break
		}
	}

}

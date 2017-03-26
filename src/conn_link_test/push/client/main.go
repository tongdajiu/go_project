package main

import (
	"conn_link_test/push/comm"
	"conn_pool"
	"fmt"
	"logger"
	"math/rand"
	"os"
	"reflect"
	"time"
	"unsafe"
)

var (
	pool       conn_pool.IConnPool
	stop_chann chan int
)

func BussParse(recv_buf []byte) ([]byte, int) {
	if len(recv_buf) >= int(unsafe.Sizeof(int16(0))) {
		packet_len := int(*(*int16)(unsafe.Pointer(&recv_buf[0])))

		if len(recv_buf) >= packet_len+int(unsafe.Sizeof(int16(0))) {
			packet_data := make([]byte, 0, packet_len+int(unsafe.Sizeof(int16(0))))
			packet_data = append(packet_data, recv_buf[0:packet_len+int(unsafe.Sizeof(int16(0)))]...)

			return packet_data, packet_len + int(unsafe.Sizeof(int16(0)))
		}
	}

	return nil, 0
}

func AfterBussConns(net_addr string, has_connected bool) (*[]byte, conn_pool.ParseFunc, interface{}) {
	var auth_req comm.AuthReq
	var slice reflect.SliceHeader
	var has_confirm bool

	has_confirm = false
	auth_req.Uuid = rand.Int63n(4000)
	copy(auth_req.User_name[:], []byte("tong"))
	auth_req.User_name[len("tong")] = 0
	copy(auth_req.Password[:], []byte("123123"))
	auth_req.Password[len("123123")] = 0

	slice.Cap = int(unsafe.Sizeof(auth_req))
	slice.Len = int(unsafe.Sizeof(auth_req))
	slice.Data = (uintptr)(unsafe.Pointer(&auth_req))

	buf := *((*[]byte)(unsafe.Pointer(&slice)))

	// 发送验证包
	message_name := "AUTH"
	send_buf := comm.Pack(buf, message_name)

	return &send_buf, BussParse, &has_confirm
}

func BussPacketHandle(connect_id int64, net_addr string, packet_buf []byte, contxt interface{}) {
	var auth_res *comm.AuthRes
	log_obj := logger.Instance()
	has_confirm := contxt.(*bool)

	buf, message_name, _ := comm.Unpack(packet_buf)
	if !(*has_confirm) {
		if message_name != "AUTH" {
			log_obj.LogAppError("Packet Not Right!Need To Be Authenticated")
			stop_chann <- 1

			return
		}
		auth_res = (*comm.AuthRes)(unsafe.Pointer(&buf[0]))
		if auth_res.Err_code == 0 {
			log_obj.LogAppInfo(" Authenticated Ok!")
			*has_confirm = true

		} else {
			log_obj.LogAppWarn("Authenticated Failed!")
			stop_chann <- 1

			return
		}
	} else {
		if message_name == "NOTICE" {
			var user_req *comm.UserPushReq
			var user_res comm.UserPushRes

			user_req = (*comm.UserPushReq)(unsafe.Pointer(&buf[0]))

			var i int
			for i = 0; user_req.Push_data[i] != 0; i++ {
			}

			log_obj.LogAppInfo("MessageName=%s,Request Data=%s", message_name, string(user_req.Push_data[:i]))

			// 应答
			user_res.Content_id = user_req.Content_id
			user_res.Stat = 1

			var slice reflect.SliceHeader

			slice.Cap = int(unsafe.Sizeof(user_res))
			slice.Len = int(unsafe.Sizeof(user_res))
			slice.Data = (uintptr)(unsafe.Pointer(&user_res))
			body := *((*[]byte)(unsafe.Pointer(&slice)))

			send_buf := comm.Pack(body, message_name)

			pool.Send2(connect_id, send_buf)
		}
	}
}

func AfterSendErr(net_addr string, packet_buf []byte) {

}

func main() {
	var conf conn_pool.ConnsCfg

	stop_chann = make(chan int)
	log_obj := logger.Instance()
	if err := log_obj.Load("./conf/conf.xml"); err != nil {
		fmt.Printf("Load Xml Failed!ErrString=%s\n", err.Error())
		return
	}

	rand.Seed(time.Now().Unix())

	//buf := comm.Pack(nil, "ALIVE")

	conf.Alram_time = 2
	conf.Conns_nums = 10
	conf.Heart_bit = nil
	conf.Max_packet_size = 1024
	conf.Set_func.Af_conn_func = AfterBussConns
	conf.Set_func.Af_send_err = AfterSendErr
	conf.Set_func.Packet_func = BussPacketHandle

	pool = conn_pool.NewConnPool("Client", &conf)
	pool.AddConns("192.168.71.13:12306")

	log_obj.LogAppInfo("Client Start!")

	for {
		<-stop_chann
		pool.Stop()

		pool := conn_pool.NewConnPool("Client", &conf)
		pool.AddConns("192.168.71.13:12306") // 重新建立连接
	}
}

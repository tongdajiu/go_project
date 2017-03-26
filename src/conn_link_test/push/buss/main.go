package main

import (
	"conn_link_test/pb"
	"conn_link_test/push/comm"
	"fmt"
	"logger"
	"math/rand"
	"reflect"
	//"runtime"
	"sync"
	"time"
	"unsafe"

	"conn_pool"

	"github.com/golang/protobuf/proto"
)

func TimeFormat(time_stamp int64) string {

	timet := time.Unix(time_stamp, 0)
	year, month, day := timet.Date()
	date_format := fmt.Sprintf("%d-%02d-%02d:%02d:%02d:%02d",
		year,
		month,
		day,
		timet.Hour(),
		timet.Minute(),
		timet.Second())

	return date_format
}

type uuidMap map[int64]string

type uuidData struct {
	uid_map   uuidMap
	rw_locker *sync.RWMutex
}

var (
	pool conn_pool.IConnPool
	//uuid_map_array []uuidData
)

func BussParse(recv_buf []byte) ([]byte, int) {
	if len(recv_buf) >= int(unsafe.Sizeof(int32(0))) {
		packet_len := int(*(*int32)(unsafe.Pointer(&recv_buf[0])))

		if len(recv_buf) >= packet_len+int(unsafe.Sizeof(int32(0))) {
			packet_data := make([]byte, 0, packet_len+int(unsafe.Sizeof(int32(0))))
			packet_data = append(packet_data, recv_buf[0:packet_len+int(unsafe.Sizeof(int32(0)))]...)

			return packet_data, packet_len + int(unsafe.Sizeof(int32(0)))
		}
	}

	return nil, 0
}

type QpsData struct {
	last_time_stamp int64
	qps_nums        int32
}

func AfterBussConns(net_addr string, has_connected bool) (*[]byte, conn_pool.ParseFunc, interface{}) {
	confirm_packet := Pack(nil, "Confirm", 0)
	var qps_data QpsData

	qps_data.last_time_stamp = time.Now().Unix()
	qps_data.qps_nums = 0
	return &confirm_packet, BussParse, &qps_data
}

func BussPacketHandle(connect_id int64, net_addr string, packet_buf []byte, contxt interface{}) {
	var user_buf []byte

	buf, message_name, _ := Unpack(packet_buf) // buss协议
	log_obj := logger.Instance()
	if message_name == "Confirm" {

		// 不做任何处理
	} else if message_name == "Pull" {

		var pull buss.Pull

		proto.Unmarshal(buf, &pull)
		user_data := pull.GetUserData()
		user_buf, message_name, _ = comm.Unpack([]byte(user_data))

		if message_name == "AUTH" {
			var auth_res comm.AuthRes
			var auth_req *comm.AuthReq
			var slice reflect.SliceHeader

			auth_req = (*comm.AuthReq)(unsafe.Pointer(&user_buf[0]))

			log_obj := logger.Instance()
			log_obj.LogAppInfo("Auth Client Info!ConnectId=%d,LinkId=%d,ThreadId=%d,UserName=%s,Password=%s",
				pull.GetConnecId(),
				pull.GetLinkId(),
				pull.GetIoThreadId(),
				string(auth_req.User_name[:]),
				string(auth_req.Password[:]))

			// 记录uuid的信息
			//index := auth_req.Uuid % len(uuid_map_array)

			//uuid_map_array[index].rw_locker.Lock()
			//uuid_map_array[index].uid_map[auth_req.Uuid] = net_addr
			//uuid_map_array[index].rw_locker.Unlock()

			auth_res.Err_code = 0
			auth_res.Err_string[0] = 0
			slice.Cap = int(unsafe.Sizeof(auth_res))
			slice.Len = int(unsafe.Sizeof(auth_res))
			slice.Data = (uintptr)(unsafe.Pointer(&auth_res))

			buf = *((*[]byte)(unsafe.Pointer(&slice)))
			user_buf = comm.Pack(buf, "AUTH")
			pull.UserData = user_buf

			var err error
			buf, err = proto.Marshal(&pull)
			if err != nil {
				log_obj.LogAppError("ProtoBuf Marshal Failed!ErrString=%s", err.Error())
			}
			res_buf := Pack(buf, "Pull", 0)
			pool.Send2(connect_id, res_buf)

			// 立即推送数据
			SendPushData(pool, connect_id, pull.GetLinkId())
		}

	} else if message_name == "Push" {
		var push_res buss.PushRes
		var user_push_res *comm.UserPushRes

		// 这里使用uuid作为link_id
		proto.Unmarshal(buf, &push_res)
		link_id := push_res.GetLinkId()
		user_data := push_res.GetUserData()
		result_code := push_res.GetResultCode()

		user_buf, message_name, _ = comm.Unpack([]byte(user_data))

		log_obj.LogAppDebug("Push Response-Data Recviced.ResultCode=%d,LinkId=%d", result_code, link_id)

		if result_code != buss.ResultCode_CLIENT_RESPONSE_OK {
			log_obj.LogAppInfo("Push Response-Data Err!ResultCode=%d,LinkId=%d", result_code, link_id)
			return
		}

		if message_name == "NOTICE" {
			user_push_res = (*comm.UserPushRes)(unsafe.Pointer(&user_buf[0]))

			log_obj.LogAppDebug("Push Data Over!Uuid=%d,ContentId=%d", link_id, user_push_res.Content_id)

			qps_data := contxt.(*QpsData)
			qps_data.qps_nums++
			if time.Now().Unix()-qps_data.last_time_stamp >= 2 { // 两秒采样
				star_time := TimeFormat(qps_data.last_time_stamp)
				end_time := TimeFormat(time.Now().Unix())
				total_time := time.Now().Unix() - qps_data.last_time_stamp

				log_obj.LogAppInfo("StartTime=%s,EndTime=%s,QpsTImes=%d", star_time, end_time, int64(qps_data.qps_nums)/total_time)
				qps_data.qps_nums = 0
				qps_data.last_time_stamp = time.Now().Unix()
			}

			// 再次推送数据
			SendPushData(pool, connect_id, link_id)
		}

	}
}

func AfterSendErr(net_addr string, buf []byte) {

}

func SendPushData(pool conn_pool.IConnPool, connect_id int64, link_id int64) {
	output_data := []string{
		"hello world",
		"jimmy",
		"tongfangwei",
	}
	content_ids := []int64{
		123232,
		532321,
		112323123,
	}

	var user_push_req comm.UserPushReq
	var slice reflect.SliceHeader
	var push_req buss.PushReq

	log_obj := logger.Instance()
	index := rand.Intn(len(content_ids))
	user_push_req.Content_id = content_ids[index]

	copy(user_push_req.Push_data[:], []byte(output_data[index]))
	user_push_req.Push_data[len(output_data[index])] = 0
	slice.Cap = int(unsafe.Sizeof(user_push_req))
	slice.Len = int(unsafe.Sizeof(user_push_req))
	slice.Data = (uintptr)(unsafe.Pointer(&user_push_req))

	buf := *((*[]byte)(unsafe.Pointer(&slice)))
	req_buf := comm.Pack(buf, "NOTICE")

	push_req.LinkId = proto.Int64(link_id)
	push_req.UserData = req_buf
	buf, _ = proto.Marshal(&push_req)
	req_buf = Pack(buf, "Push", 0)

	log_obj.LogAppDebug("Push Data Again!ContentId=%d,LinkId=%d,Data=%s",
		content_ids[index],
		link_id,
		output_data[index])

	pool.Send2(connect_id, req_buf)

}

func Initialize() {
	var funcs conn_pool.ConnFunc

	log_obj := logger.Instance()
	err := log_obj.Load("./conf/conf.xml")
	if err != nil {
		fmt.Printf("Load Failed!ErrString=%s", err.Error())
		return
	}

	//array_size := runtime.NumCPU()
	//uuid_map_array = make([]uuidData, array_size)

	/*
		for i := 0; i < len(uuid_map_array); i++ {
			uuid_map_array[i] = uuidData{
				uid_map:   make(uuidMap),
				rw_locker: &sync.RWMutex{},
			}
		}

	*/
	funcs.Af_conn_func = AfterBussConns
	funcs.Packet_func = BussPacketHandle
	funcs.Af_send_err = AfterSendErr

	// 设置心跳包
	//heart_bit := Pack(nil, "Ping", 0)

	conf := conn_pool.ConnsCfg{
		Max_packet_size: 1024,
		Heart_bit:       nil,
		Alram_time:      5,
		Conns_nums:      5,
		Set_func:        funcs,
	}

	pool = conn_pool.NewConnPool("buss", &conf)
	pool.AddConns("192.168.71.13:12307")

	log_obj.LogAppInfo("Buss Init Ok!")
}

func main() {
	// 初始化
	Initialize()

	for {
		time.Sleep(10 * time.Second)
	}
}

package main

import (
	"conn_pool"
	"container/list"
	"logger"
	"unsafe"
	//"websocket_test/pb"
	"game_core"

	"github.com/gorilla/websocket"
	//"github.com/golang/protobuf/proto"
)

type linkData struct {
	pool             conn_pool.IConnPool
	remote_ip        string
	bk_server_list   *list.List
	bk_has_connected chan bool
	web_client       *WebClient
}

func AfConFunc(net_addr string, has_connected bool, contxt interface{}) (*[]byte, conn_pool.ParseFunc) {
	link_data := contxt.(*linkData)
	if !has_connected {

		// 没有连接成功
		pool := link_data.pool
		pool.RmConns(link_data.remote_ip)
		if link_data.bk_server_list.Len() <= 0 {
			link_data.bk_has_connected <- false
		}

		remote_ip := link_data.bk_server_list.Remove(link_data.bk_server_list.Front()).(string)
		pool.AddConns(remote_ip)
	} else {
		link_data.bk_has_connected <- true // 连接成功
	}

	return nil, game_core.GameParse
}

func PacketFunc(net_addr string, packet_buf []byte, contxt interface{}) {

}

func AfSendErr(net_addr string, buf []byte, contxt interface{}) {
	// 连接断开,内部服务器
	web_client := contxt.(*WebClient)
	web_client.Close()
}

type webLogic struct {
	forward_svr []string
}

func (this *webLogic) OnPacketRecv(wc *WebClient, buf []byte, contxt interface{}) {
	log_obj := logger.Instance()
	pool_data := contxt.(*linkData)
	conn_pool := pool_data.pool
	remote_ip := pool_data.remote_ip
	packet_header := (*game_core.PacketHeader)(unsafe.Pointer(&buf[0]))

	check_sum := game_core.CheckSum(buf[4:])
	if check_sum != packet_header.Check_sum {
		log_obj.LogAppWarn("CheckSum Failed!calc_checksum=%d,checksum=%d,PeerIP=%s",
			check_sum,
			packet_header.Check_sum,
			wc.GetPeerAddr())
		wc.Close()
	}

	// 数据包转发
	conn_pool.Send(remote_ip, buf)

}

func (this *webLogic) OnErr(wc *WebClient, contxt interface{}) {
	// 断开连接
	link_data := contxt.(*linkData)
	link_data.pool.Stop()
}

func (this *webLogic) LinkWebSocket(ws *websocket.Conn) *WebClient {
	// 创建后台连接一一对应
	set_func := conn_pool.ConnFunc{
		Af_conn_func: AfConFunc,
		Packet_func:  PacketFunc,
		Af_send_err:  AfSendErr,
	}

	conf := conn_pool.ConnsCfg{
		1024,
		nil,
		0,
		1,
		set_func}

	link_data := &linkData{}
	link_data.bk_server_list = list.New()

	link_data.pool = conn_pool.NewConnPool("bk_sever", &conf, link_data)
	link_data.remote_ip = "127.0.0.1:12345"
	link_data.bk_has_connected = make(chan bool, 1)
	link_data.pool.AddConns("127.0.0.1:12345")

	var has_connected bool
	select {
	case has_connected = <-link_data.bk_has_connected:
		break
	}

	if !has_connected {
		return nil
	}
	// 连接成功,绑定前端和后端的连接
	web_client := NewWebClient(ws,
		svr_config.max_packet_size,
		1000,
		svr_config.alram_time,
		link_data,
		this,
		game_core.GameParse)

	link_data.web_client = web_client

	return web_client
}

package main

import (
	"logger"
	"time"

	"github.com/gorilla/websocket"
)

type IWebListener interface {
	OnPacketRecv(*WebClient, []byte, interface{})
	OnErr(*WebClient, interface{})
}

type WebClient struct {
	max_packet_size int
	send_chann      chan []byte
	close_chann     chan int // 自动关掉”标志“
	web_socket      *websocket.Conn
	alarm_time      int
	shake_hand      bool
	contxt          interface{}
	web_listener    IWebListener
	parse_func      func(buf []byte) ([]byte, int)
}

func (this *WebClient) Close() {

	// "断开"连接
	this.close_chann <- 1
}

func (this *WebClient) GetPeerAddr() string {
	return this.web_socket.RemoteAddr().String()
}

func (this *WebClient) SendMessage(send_buf []byte) bool {
	time_chann := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		time_chann <- 1
	}()

	select {
	case this.send_chann <- send_buf:
		return true
	case <-time_chann:
		return false
	}
}

func NewWebClient(ws *websocket.Conn,
	max_packet_size int,
	max_send_packet_nums int,
	alarm_time int,
	contxt interface{},
	web_listener IWebListener,
	parse_func func(buf []byte) ([]byte, int)) *WebClient {

	var web_client *WebClient = &WebClient{
		web_socket:      ws,
		send_chann:      make(chan []byte, max_send_packet_nums),
		close_chann:     make(chan int),
		max_packet_size: max_packet_size,
		alarm_time:      alarm_time,
		contxt:          contxt,
		web_listener:    web_listener,
		parse_func:      parse_func,
	}

	go web_client.sendRounte()
	return web_client
}

func (this *WebClient) DoReadMessage() {
	pack_buf := make([]byte, this.max_packet_size)
	parse_offset := 0
	end := 0
	log_obj := logger.Instance()

	defer func() {
		this.web_listener.OnErr(this, this.contxt)
		this.web_socket.Close()
	}()
	for {
		dead_line := time.Now().Add(time.Duration(this.alarm_time) * time.Second)
		this.web_socket.SetReadDeadline(dead_line)
		_, data, err := this.web_socket.ReadMessage()
		if err != nil {
			log_obj.LogSysError("ReadMessage Failed!RemoteAddr=%s,ErrString=%s",
				this.web_socket.RemoteAddr(),
				err.Error())
			return
		}

		// 空间不够
		if this.max_packet_size-end <= len(data) {

			// 重定义parse_offset, end
			copy(pack_buf[0:], pack_buf[parse_offset:end])
			end -= parse_offset
			parse_offset = 0

			if this.max_packet_size-end <= len(data) {
				log_obj.LogSysError("Packet Buffer Full!RemoteAddr=%s", this.web_socket.RemoteAddr())
				return
			}
		}
		copy(pack_buf[end:], data)
		end += len(data)

		for {
			buf, size_parsed := this.parse_func(pack_buf[parse_offset:end])
			if size_parsed == 0 {
				break
			}

			// 收到一个完整的数据包
			this.web_listener.OnPacketRecv(this, buf, this.contxt)
			parse_offset += size_parsed
		}
	}
}

func (this *WebClient) sendRounte() {
	defer this.web_socket.Close()
	log_obj := logger.Instance()
	for {
		select {
		case send_buf := <-this.send_chann:
			err := this.web_socket.WriteMessage(websocket.BinaryMessage, send_buf)
			if err != nil {
				log_obj.LogSysError("WriteMessage Failed!RemoteAddr=%s,ErrString=%s",
					this.web_socket.RemoteAddr(),
					err.Error())
				return
			}
		case <-this.close_chann:
			log_obj.LogSysWarn("Close WebClient!RemoteAddr=%s", this.web_socket.RemoteAddr())
			return
		}
	}
}

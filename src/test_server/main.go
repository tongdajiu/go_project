// server project main.go
package main

import (
	"fmt"
	"logger"
	"net_server"
	"time"
)

//go build -gcflags "-N -l"
//go -o websocket_test build -gcflags "-N -l"

//info goroutines
type ClientLogic struct {
	log_obj logger.ILogger
}

func (this *ClientLogic) OnNetParsing(client net_server.INetClient,
	rev_buf []byte,
	contxt interface{}) int32 {
	var parsed_size int32

	parsed_size = 0
	packet_size := len(string("hello world"))
	if len(rev_buf) >= packet_size {
		log := logger.Instance()
		data := string(rev_buf)
		log.LogAppDebug("recv_data=%s", data)

		parsed_size = int32(packet_size)
		client.Send(rev_buf[:parsed_size])
	}

	return parsed_size
}

func (this *ClientLogic) OnNetAccpet(peer_ip string) (interface{}, bool) {
	return nil, true
}

func (this *ClientLogic) OnNetErr(client net_server.INetClient, contxt interface{}) {

}

func main() {
	var config net_server.Config
	var client_logic ClientLogic

	log_obj := logger.Instance()
	err := log_obj.Load("./conf.xml")
	if err != nil {
		fmt.Println("Load Log Conf Failed!ErrString=%s", err.Error())
		return
	}

	defer log_obj.Close()
	config.SetServerAddr("0:0:0:0", "50123")
	config.SetPacketSize(1024)
	config.SetAccpetNums(1000)
	config.SetOutPacketNums(4000)
	config.SetAlramTime(10)

	server := net_server.CreateNetServer(&client_logic)
	if !server.Start(&config) {
		log_obj.LogAppError("Bind Failed!")
		return
	}

	time.Sleep(1 * time.Hour)
}

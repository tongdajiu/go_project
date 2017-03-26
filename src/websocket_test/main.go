package main

import (
	"container/list"
	"flag"
	"fmt"
	"logger"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

func MyCheckOrigin(rsp *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     MyCheckOrigin,
}

type SvrConfig struct {
	alram_time      int
	max_packet_size int
	bk_server       list.List
}

var (
	addr       = flag.String("addr", ":8080", "http service address")
	homeTempl  *template.Template
	web_logic  webLogic
	svr_config = SvrConfig{
		alram_time:      20,
		max_packet_size: 1024,
	}
)

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func wsHandler(w http.ResponseWriter, req *http.Request) {
	log_obj := logger.Instance()
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log_obj.LogAppWarn("shake hanlde Failed!ErrString=%s", err.Error())
		return
	}

	defer ws.Close()
	web_client := web_logic.LinkWebSocket(ws)
	if web_client == nil {
		log_obj.LogAppWarn("Link Faild!")
		return
	}

	web_client.DoReadMessage()
}

type AA struct {
	a int
	b int
}

func main() {

	log_obj := logger.Instance()
	homeTempl = template.Must(template.ParseFiles("../html/login.html"))
	if err := log_obj.Load("../conf/log.xml"); err != nil {
		fmt.Printf("Logger Load Failed!ErrString=%s", err.Error())
		return
	}

	log_obj.LogAppDebug("ListenAndServe:Start!")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", wsHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log_obj.LogAppError("ListenAndServe:ErrString=%s", err.Error())
	}

	log_obj.LogAppDebug("ListenAndServe:Finished!")
}

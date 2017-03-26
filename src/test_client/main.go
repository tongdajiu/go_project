package main

import (
	"fmt"
	"logger"
	"net"
)

func Run(conn net.Conn, stop_chann chan int) {
	buff := make([]byte, 512)
	log := logger.Instance()

	defer conn.Close()
	for {
		buff = []byte("hello world")
		_, err := conn.Write(buff)
		if err != nil {
			log.LogAppWarn("Send Failed!ErrString=%s", err.Error())
			break
		}

		log.LogAppDebug("Send Data=%s", string(buff))

		buff[0] = 0
		_, err = conn.Read(buff)
		if err != nil {
			log.LogAppWarn("Recv Failed!ErrString=%s", err.Error())
			break
		}

		recv_data := string(buff)
		log.LogAppDebug("Recv-Data=%s", recv_data)
	}

	stop_chann <- 1
}

func main() {
	log := logger.Instance()
	var conn net.Conn

	err := log.Load("./conf.xml")
	if err != nil {
		fmt.Println("Load Xml Failed!")
		return
	}

	defer log.Close()
	conn, err = net.Dial("tcp4", "127.0.0.1:50123")
	if err != nil {
		log.LogAppError("Connect Server Failed!ErrString=%s", err.Error())
		return
	}

	stop_chann := make(chan int, 1)
	go Run(conn, stop_chann)
	select {
	case <-stop_chann:
	}
}

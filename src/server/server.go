package server

import (
	"log"
	"net"
	"syscall"
)

// 服务端
//
func Server(hostip, port string) {
	log.Println("server listening is starting")
	listener, err := net.Listen("tcp", hostip+":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		//循环接收客户端连接
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go send(conn)
		go receive(conn)
	}
}

func send(conn net.Conn) {
}

func receive(conn net.Conn) {
	// 完整数据包
	singlePackData := make([]byte, 1024)
	for {
		buf := make([]byte, 1024)
		len, err := conn.Read(buf)
		//处理粘包
		switch err {
		// 缓存区读取完
		case nil:
			// read 1024bytes
			// get data package size
			buf = nil

		case syscall.EAGAIN:
			continue
		default:
			goto DISCONNECT
		}
		//获取序列
		ch := make(chan []byte, 1)
		chatContainer["1"] = ch

	}
DISCONNECT:
	err := conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}

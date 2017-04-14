package main

import (
	"github.com/weihualiu/chat-demo/src/server"
	"log"
	//"net"
)

func main() {
	log.Println("chat demo start")
	server.Server("127.0.0.1", "50000")
}

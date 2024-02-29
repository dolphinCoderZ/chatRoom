package main

import (
	"chatRoom/server/model"
	"fmt"
	"net"
	"time"
)

func handle(conn net.Conn) {
	defer conn.Close()

	processor := &Processor{
		Conn: conn,
	}
	processor.ServerProcess()
}

func initUserDAO() {
	model.GlobalUserDAO = model.NewUserDAO(pool)
}

func init() {
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDAO()
}

func main() {
	fmt.Println("服务器监听8889")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("Listen err=", err)
		return
	}

	for {
		fmt.Println("等待客户端连接")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("server accept err=", err)
		}
		go handle(conn)
	}
}

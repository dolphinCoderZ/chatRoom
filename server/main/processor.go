package main

import (
	"chatRoom/common/message"
	"chatRoom/server/process"
	"chatRoom/server/utils"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

func serverProcessByMesType(conn net.Conn, mesRead *message.Message) (processErr error) {
	switch mesRead.Type {
	case message.LoginMesType:
		userProcessor := &process.UserProcessor{
			Conn: conn,
		}
		processErr = userProcessor.ServerProcessLogin(mesRead)
	case message.RegisterMesType:
		userProcessor := &process.UserProcessor{
			Conn: conn,
		}
		processErr = userProcessor.ServerProcessRegister(mesRead)
	case message.LogOutMesType:
		userProcessor := &process.UserProcessor{
			Conn: conn,
		}
		processErr = userProcessor.ServerProcessLogOut(mesRead)
	case message.SmsMesType:
		smsProcessor := &process.SmsProcessor{}
		processErr = smsProcessor.SendMesToOnlineUsers(mesRead)
	case message.OneToOneMesType:
		smsProcessor := &process.SmsProcessor{}
		processErr = smsProcessor.SendMesByUserId(mesRead)
	default:
		fmt.Println("消息类型错误，无法处理")
	}
	return processErr
}

func (this *Processor) ServerProcess() {
	for {
		transfer := &utils.Transfer{
			Conn: this.Conn,
		}
		//监听客户端发来的信息
		mesRead, err := transfer.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端关闭")
				return
			} else {
				fmt.Println("ServerProcess ReadPkg err")
				return
			}
		}
		//fmt.Println(mesRead)

		err = serverProcessByMesType(this.Conn, &mesRead)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

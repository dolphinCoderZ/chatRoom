package process

import (
	"chatRoom/client/utils"
	"chatRoom/common/message"
	"fmt"
	"net"
	"os"
)

func showMenu() {
	var choice int
	var content string
	//进入二级菜单，就有可能群发消息或者私聊
	smsProcess := &SmsProcess{}
	for {
		fmt.Println("-----------恭喜xxx登录成功-----------")
		fmt.Println("\t 1.显示在线用户列表")
		fmt.Println("\t 2.群发消息")
		fmt.Println("\t 3.私聊")
		fmt.Println("\t 4.退出系统")
		fmt.Println("请选择（1-4）:")

		_, _ = fmt.Scanf("%d\n", &choice)
		switch choice {
		case 1:
			showOnlineUsers()
		case 2:
			fmt.Println("您想对大家说点什么:")
			_, _ = fmt.Scanf("%s\n", &content)
			if content == "" {
				fmt.Println("消息不能为空")
				continue
			}
			sendErr := smsProcess.SendGroupMes(content)
			if sendErr != nil {
				fmt.Println("发送失败，请重新发送")
				fmt.Println()
			} else {
				fmt.Println("恭喜您，发送成功")
				fmt.Println()
			}
		case 3:
			var ReceiveUserId int
			fmt.Println("请选择你想私聊的用户ID：")
			_, _ = fmt.Scanf("%d\n", &ReceiveUserId)
			_, ok := userManager.onlineUsers[ReceiveUserId]
			if !ok {
				fmt.Println("该用户当前不在线，请重新选择")
				continue
			}
			fmt.Printf("您想对用户%d说点什么：", ReceiveUserId)
			_, _ = fmt.Scanf("%s\n", &content)
			if content == "" {
				fmt.Println("消息不能为空")
				continue
			}
			sendErr := smsProcess.SendOneToOneMes(ReceiveUserId, content)
			if sendErr != nil {
				fmt.Println("发送失败，请重新发送")
				fmt.Println()
			} else {
				fmt.Println("恭喜您，发送成功")
				fmt.Println()
			}
		case 4:
			userProcess := &UserProcess{}
			userProcess.LogOut()
			os.Exit(0)
		default:
			fmt.Println("请重新选择")
		}
	}
}

func listenServerMes(conn net.Conn) {
	transfer := &utils.Transfer{
		Conn: conn,
	}
	for {
		listenMes, err := transfer.ReadPkg()
		if err != nil {
			fmt.Println("listen message err", err)
			return
		}
		switch listenMes.Type {
		case message.NotifyUserStatusMesType:
			updateUserStatus(&listenMes)
		case message.SmsMesType:
			smsManager.showGroupMes(&listenMes)
		case message.OneToOneResMesType:
			smsManager.showOneToOneMes(&listenMes)
		case message.LogOutMesType:
			deleteLogOutUsers(&listenMes)
		default:
			fmt.Println("无法处理该消息")
		}
	}
}

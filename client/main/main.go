package main

import (
	"chatRoom/client/process"
	"fmt"
	"os"
)

// 用户id和密码
var userId int
var password string
var userName string

func main() {
	var choice int

	for {
		fmt.Println("-------------欢迎登录多人聊天系统-------------")
		fmt.Println("\t\t 1 登录聊天室")
		fmt.Println("\t\t 2 注册用户")
		fmt.Println("\t\t 3 退出系统")
		fmt.Println("\t\t 请重新选择（1-3）:")

		_, _ = fmt.Scanf("%d\n", &choice)
		switch choice {
		case 1:
			fmt.Println("请输入用户ID:")
			//fmt.Scanln(&userId)
			_, _ = fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入密码:")
			_, _ = fmt.Scanf("%s\n", &password)

			userProcess := &process.UserProcess{}
			err := userProcess.Login(userId, password)
			if err != nil {
				fmt.Println(err)
				return
			}
		case 2:
			fmt.Println("注册用户")
			for {
				fmt.Println("请输入用户ID:")
				userIdN, _ := fmt.Scanf("%d\n", &userId)
				fmt.Println("请输入密码:")
				passwordN, _ := fmt.Scanf("%s\n", &password)
				fmt.Println("请输入用户昵称:")
				_, _ = fmt.Scanf("%s\n", &userName)
				if userIdN != 0 && passwordN != 0 {
					break
				}
				fmt.Println("用户ID和密码都不能为空，请重新注册")
				fmt.Println()
			}

			userProcess := &process.UserProcess{}
			err := userProcess.Register(userId, password, userName)
			if err != nil {
				fmt.Println(err)
			} else {
				continue
			}
		case 3:
			os.Exit(0)
		default:
			fmt.Println("输入有误，请重新输入")
		}
	}
}

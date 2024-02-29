package process

import (
	"chatRoom/client/utils"
	"chatRoom/common/message"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type UserProcess struct {
}

//注册和登录在一级菜单，都需要拨号连接服务器

func (this *UserProcess) Register(userId int, password string, userName string) (err error) {
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		err = errors.New("register net.Dial err")
		return
	}
	defer conn.Close()

	var mes message.Message
	mes.Type = message.RegisterMesType

	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = password
	registerMes.User.UserName = userName

	registerData, err := json.Marshal(registerMes.User)
	if err != nil {
		err = errors.New("register json.Marshal(registerMes.User) err")
		return
	}
	mes.Data = string(registerData)
	registerMesData, err := json.Marshal(mes)
	if err != nil {
		err = errors.New("register json.Marshal(mes) err")
		return
	}

	//发送注册信息给服务器
	transfer := utils.Transfer{
		Conn: conn,
	}
	err = transfer.WritePkg(registerMesData)
	if err != nil {
		err = errors.New("register WritePkg err")
		return
	}

	//服务器返回的注册校验信息
	registerRes, err := transfer.ReadPkg()
	if err != nil {
		err = errors.New("register ReadPkg err")
		return
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(registerRes.Data), &registerResMes)
	if err != nil {
		err = errors.New("register json.Unmarshal(registerRes.Data) err")
	}

	if registerResMes.Code == 200 {
		fmt.Println("注册成功")
	} else {
		err = errors.New(registerResMes.Error)
		return
	}
	return
}

func (this *UserProcess) Login(userId int, password string) (err error) {
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		err = errors.New("login net.Dial err")
		return
	}
	defer conn.Close()

	var mes message.Message
	mes.Type = message.LoginMesType

	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = password
	loginData, err := json.Marshal(loginMes)
	if err != nil {
		err = errors.New("login json.Marshal(loginMes) err")
		return
	}
	mes.Data = string(loginData)
	loginMesData, err := json.Marshal(mes)
	if err != nil {
		err = errors.New("login json.Marshal(mes) err")
	}

	transfer := &utils.Transfer{
		Conn: conn,
	}
	err = transfer.WritePkg(loginMesData)
	if err != nil {
		err = errors.New("login WritePkg err")
		return
	}

	//服务器返回的登录验证消息
	loginRes, err := transfer.ReadPkg()
	if err != nil {
		err = errors.New("login ReadPkg err")
		return
	}

	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(loginRes.Data), &loginResMes)
	if err != nil {
		err = errors.New("login json.Unmarshal(loginRes.Data) err")
		return
	}

	if loginResMes.Code == 200 {
		currentUser.Conn = conn
		currentUser.UserId = userId
		currentUser.UserStatus = message.UserOnline
		//某个用户登录成功，先从登录返回信息中获取目前所有在线的用户
		for _, onlineId := range loginResMes.OnlineUsersId {
			user := &message.User{
				UserId:     onlineId,
				UserStatus: message.UserOnline,
			}
			//初始在线用户列表
			userManager.onlineUsers[onlineId] = user
		}
		showOnlineUsers()
		//每个conn(每个客户端)，保持与服务器通讯，监控服务器的信息
		//服务器发送信息，客户端根据信息的类别进行相应的处理
		//1.上线通知消息：更新客户端在线列表，并显示最新在线用户列表
		//2.下线通知消息：更新客户端在线列表，并显示最新在线用户列表
		//3.群发消息通知：展示群发消息
		//4.私聊消息通知：展示私聊消息
		go listenServerMes(conn)
		showMenu()
	} else {
		err = errors.New(loginResMes.Error)
		return
	}
	return
}

// 问题：更新userManager.onlineUsers需要由服务器通知其他用户
func (this *UserProcess) LogOut() {
	var mes message.Message
	mes.Type = message.LogOutMesType

	var logOutMes message.LogOutMes
	logOutMes.UserId = currentUser.UserId
	logOutInfo, err := json.Marshal(logOutMes)
	if err != nil {
		fmt.Println("logout json.Marshal(logOutMes) err")
		return
	}
	mes.Data = string(logOutInfo)

	resMes, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("logout json.Marshal(mes) err")
		return
	}

	transfer := &utils.Transfer{
		Conn: currentUser.Conn,
	}
	err = transfer.WritePkg(resMes)
	if err != nil {
		fmt.Println("logOut WritePkg err")
	}
}

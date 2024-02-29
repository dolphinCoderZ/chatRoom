package process

import (
	"chatRoom/common/message"
	"chatRoom/server/model"
	"chatRoom/server/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

type UserProcessor struct {
	Conn   net.Conn
	UserId int
}

func Notify(conn net.Conn, userId int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	notifyData, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("Notify json.Marshal err")
		return
	}
	mes.Data = string(notifyData)

	resMes, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("Notify json.Marshal err")
		return
	}

	transfer := &utils.Transfer{
		Conn: conn,
	}
	err = transfer.WritePkg(resMes)
	if err != nil {
		fmt.Println("Notify WritePkg err")
		return
	}
}

func (this *UserProcessor) NotifyOtherOnlineUsers() {
	for id, userProcess := range userManager.onlineUsers {
		if id == this.UserId {
			continue
		}
		//遍历通知其他每个在线客户端
		Notify(userProcess.Conn, this.UserId)
	}
}

func (this *UserProcessor) ServerProcessLogin(mesRead *message.Message) (loginErr error) {
	var loginMes message.LoginMes
	loginErr = json.Unmarshal([]byte(mesRead.Data), &loginMes)
	if loginErr != nil {
		loginErr = errors.New("ServerProcessLogin json.Unmarshal err")
		return
	}

	var resMes message.Message
	resMes.Type = message.LoginResMesType

	var loginResMes message.LoginResMes
	//由于客户端是在登录成功之后才拿到在线用户，未成功之前都无法取到在线用户，因此客户端无法校验该用户是否已经登录
	_, ok := userManager.onlineUsers[loginMes.UserId]
	if ok {
		loginResMes.Code = 402
		loginResMes.Error = model.ErrorAlreadyLogin.Error()
	} else {
		//数据库校验
		user, loginErr := model.GlobalUserDAO.LoginCheck(loginMes.UserId, loginMes.UserPwd)
		if loginErr != nil {
			if errors.Is(loginErr, model.ErrorUserNotExists) {
				loginResMes.Code = 500
				loginResMes.Error = loginErr.Error()
			} else if errors.Is(loginErr, model.ErrorUserPwd) {
				loginResMes.Code = 403
				loginResMes.Error = loginErr.Error()
			} else {
				loginResMes.Code = 404
				loginResMes.Error = "未知数据库错误"
			}
		} else {
			loginResMes.Code = 200
			this.UserId = loginMes.UserId
			userManager.AddOnlineUser(this)
			//某个客户端上线，通知其他在线客户端："我已上线"
			this.NotifyOtherOnlineUsers()
			//返回给登录成功客户端当前在线用户ID列表
			for id, _ := range userManager.onlineUsers {
				loginResMes.OnlineUsersId = append(loginResMes.OnlineUsersId, id)
			}
			fmt.Println(user, "登录成功")
		}
	}

	loginResData, loginErr := json.Marshal(loginResMes)
	if loginErr != nil {
		loginErr = errors.New("ServerProcessLogin json.Marshal(loginResMes) err")
		return
	}

	resMes.Data = string(loginResData)
	resData, loginErr := json.Marshal(resMes)
	if loginErr != nil {
		loginErr = errors.New("ServerProcessLogin json.Marshal(resMes) err")
		return
	}

	transfer := &utils.Transfer{
		Conn: this.Conn,
	}
	loginErr = transfer.WritePkg(resData)
	if loginErr != nil {
		loginErr = errors.New("ServerProcessLogin WritePkg err")
	}
	return
}

func (this *UserProcessor) ServerProcessRegister(mesRead *message.Message) (registerErr error) {
	var registerMes message.RegisterMes
	//mesRead.Data反序列化直接对应registerMes.User
	registerErr = json.Unmarshal([]byte(mesRead.Data), &registerMes.User)
	if registerErr != nil {
		registerErr = errors.New("ServerProcessRegister json.Unmarshal err")
		return
	}
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//注册数据库
	registerErr = model.GlobalUserDAO.Register(&registerMes.User)
	if registerErr != nil {
		if errors.Is(registerErr, model.ErrorExists) {
			registerResMes.Code = 400
			registerResMes.Error = registerErr.Error()
		} else {
			registerResMes.Code = 505
			registerResMes.Error = "注册出现未知错误"
		}
	} else {
		registerResMes.Code = 200
	}

	registerResData, registerErr := json.Marshal(registerResMes)
	if registerErr != nil {
		registerErr = errors.New("ServerProcessRegister json.Marshal(registerResMes) err")
		return
	}
	resMes.Data = string(registerResData)
	resData, registerErr := json.Marshal(resMes)
	if registerErr != nil {
		registerErr = errors.New("ServerProcessRegister json.Marshal(resMes) err")
		return
	}

	transfer := utils.Transfer{
		Conn: this.Conn,
	}
	registerErr = transfer.WritePkg(resData)
	if registerErr != nil {
		registerErr = errors.New("ServerProcessRegister WritePkg err")
		return
	}
	return
}

func (this *UserProcessor) ServerProcessLogOut(mesRead *message.Message) (logOutErr error) {
	var logOutMes message.LogOutMes
	logOutErr = json.Unmarshal([]byte(mesRead.Data), &logOutMes)
	if logOutErr != nil {
		logOutErr = errors.New("ServerProcessLogOut json.Unmarshal err")
		return
	}
	logOutUserId := logOutMes.UserId
	userManager.DeleteOnlineUser(logOutUserId)
	fmt.Println("用户:", logOutUserId, "退出登录")

	//通知其他客户端要退出的客户端的ID等信息(退出客户端发来的原信息包就包含了所需信息)
	resMes, err := json.Marshal(mesRead)
	if err != nil {
		fmt.Println("ServerProcessLogOut json.Marshal err")
		return
	}

	//通知其他客户端
	for _, userProcess := range userManager.onlineUsers {
		transfer := &utils.Transfer{
			Conn: userProcess.Conn,
		}
		err = transfer.WritePkg(resMes)
		if err != nil {
			fmt.Println("ServerProcessLogOut WritePkg err")
			return
		}
	}
	return
}

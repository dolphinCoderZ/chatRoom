package process

import (
	"chatRoom/client/model"
	"chatRoom/common/message"
	"encoding/json"
	"fmt"
)

type UserManager struct {
	onlineUsers map[int]*message.User
}

var userManager *UserManager
var currentUser *model.CurrentUser

func init() {
	userManager = &UserManager{
		onlineUsers: make(map[int]*message.User, 10),
	}
	currentUser = &model.CurrentUser{}
}

func showOnlineUsers() {
	fmt.Println()
	fmt.Println("当前在线用户ID列表:")
	for id, _ := range userManager.onlineUsers {
		fmt.Println("用户ID:\t", id)
	}
	fmt.Println()
}

func updateUserStatus(listenMes *message.Message) {
	var notifyMes message.NotifyUserStatusMes
	err := json.Unmarshal([]byte(listenMes.Data), &notifyMes)
	if err != nil {
		fmt.Println("listenServerMes notifyMes json.Unmarshal")
	}
	user, ok := userManager.onlineUsers[notifyMes.UserId]
	if !ok {
		user = &message.User{
			UserId: notifyMes.UserId,
		}
	}
	user.UserStatus = notifyMes.Status
	//更新客户端的在线列表
	userManager.onlineUsers[notifyMes.UserId] = user
	showOnlineUsers()
}

func deleteLogOutUsers(listenMes *message.Message) {
	var logOutMes message.LogOutMes
	err := json.Unmarshal([]byte(listenMes.Data), &logOutMes)
	if err != nil {
		fmt.Println("notifyMes json.Unmarshal")
	}
	delete(userManager.onlineUsers, logOutMes.UserId)
	showOnlineUsers()
}

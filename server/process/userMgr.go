package process

import (
	"fmt"
)

// 每个客户端登录成功后，conn需要保存起来，方便服务器转发信息
type UserMgr struct {
	onlineUsers map[int]*UserProcessor
}

var userManager *UserMgr

func init() {
	userManager = &UserMgr{
		onlineUsers: make(map[int]*UserProcessor, 1024),
	}
}

// map在传入相同的key相当于更新操作
func (this *UserMgr) AddOnlineUser(userProcessor *UserProcessor) {
	userManager.onlineUsers[userProcessor.UserId] = userProcessor
}

func (this *UserMgr) DeleteOnlineUser(userId int) {
	delete(userManager.onlineUsers, userId)
}

func (this *UserMgr) GetOnlineUserById(userId int) (userProcessor *UserProcessor, err error) {
	userProcessor, ok := userManager.onlineUsers[userId]
	if !ok {
		err = fmt.Errorf("用户%d不存在", userId)
		return
	}
	return
}

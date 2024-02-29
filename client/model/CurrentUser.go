package model

import (
	"chatRoom/common/message"
	"net"
)

// 客户端发送消息时，需要用到当前用户和服务器的连接
type CurrentUser struct {
	Conn net.Conn
	message.User
}

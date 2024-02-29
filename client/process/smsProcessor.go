package process

import (
	"chatRoom/client/utils"
	"chatRoom/common/message"
	"encoding/json"
	"errors"
)

type SmsProcess struct {
}

//不论群发私发，都是发送给服务器，由服务器解析消息包转发给对应的客户端

func (this *SmsProcess) SendGroupMes(content string) (sendMesErr error) {
	var mes message.Message
	mes.Type = message.SmsMesType
	var contentMes message.SmsMes

	contentMes.Content = content
	contentMes.UserId = currentUser.UserId
	contentMes.UserStatus = currentUser.UserStatus
	contentData, sendMesErr := json.Marshal(contentMes)
	if sendMesErr != nil {
		sendMesErr = errors.New("SendGroupMes json.Marshal(contentMes) err")
		return
	}
	mes.Data = string(contentData)

	smsData, sendMesErr := json.Marshal(mes)
	if sendMesErr != nil {
		sendMesErr = errors.New("SendGroupMes json.Marshal(mes) err")
		return
	}

	transfer := &utils.Transfer{
		Conn: currentUser.Conn,
	}
	sendMesErr = transfer.WritePkg(smsData)
	if sendMesErr != nil {
		sendMesErr = errors.New("SendGroupMes WritePkg err")
		return
	}
	return
}

func (this *SmsProcess) SendOneToOneMes(ReceiveUserId int, content string) (sendMesErr error) {
	var mes message.Message
	mes.Type = message.OneToOneMesType

	var contentMes message.OneToOneMes
	contentMes.Content = content
	contentMes.ReceiveUserId = ReceiveUserId
	contentMes.UserId = currentUser.UserId
	contentMes.UserStatus = currentUser.UserStatus
	contentData, sendMesErr := json.Marshal(contentMes)
	if sendMesErr != nil {
		sendMesErr = errors.New("SendOneToOneMes json.Marshal(contentMes) err")
		return
	}
	mes.Data = string(contentData)
	oneToOneData, sendMesErr := json.Marshal(mes)
	if sendMesErr != nil {
		sendMesErr = errors.New("SendOneToOneMes json.Marshal(mes) err")
		return
	}

	transfer := &utils.Transfer{
		Conn: currentUser.Conn,
	}
	sendMesErr = transfer.WritePkg(oneToOneData)
	if sendMesErr != nil {
		sendMesErr = errors.New("SendOneToOneMes WritePkg err")
		return
	}
	return
}

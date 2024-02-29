package process

import (
	"chatRoom/client/utils"
	"chatRoom/common/message"
	"encoding/json"
	"errors"
	"net"
)

type SmsProcessor struct {
}

func (this *SmsProcessor) SendMesToOnlineUsers(mesRead *message.Message) (sendErr error) {
	//不做任何处理，直接把包转发给在线的所有客户端
	mes, sendErr := json.Marshal(mesRead)
	if sendErr != nil {
		sendErr = errors.New("SendMesToOnlineUsers json.Marshal err")
	}

	for _, userProcess := range userManager.onlineUsers {
		sendErr = sendMes(mes, userProcess.Conn)
		if sendErr != nil {
			sendErr = errors.New("SendMesToOnlineUsers sendMes err")
		}
	}
	return
}

func sendMes(data []byte, conn net.Conn) (sendErr error) {
	transfer := &utils.Transfer{
		Conn: conn,
	}
	sendErr = transfer.WritePkg(data)
	if sendErr != nil {
		sendErr = errors.New("sendMes WritePkg err")
		return
	}
	return
}

func (this *SmsProcessor) SendMesByUserId(mesRead *message.Message) (sendErr error) {
	//从发送方的信息包中获取接收方信息
	var oneToOneMes message.OneToOneMes
	sendErr = json.Unmarshal([]byte(mesRead.Data), &oneToOneMes)
	if sendErr != nil {
		sendErr = errors.New("sendMesByUserId json.Unmarshal err")
	}

	sendId := oneToOneMes.UserId
	content := oneToOneMes.Content
	receiveId := oneToOneMes.ReceiveUserId
	receiveUserProcess := userManager.onlineUsers[receiveId]

	var mes message.Message
	mes.Type = message.OneToOneResMesType

	var ontToOneResMes message.OneToOneResMes
	ontToOneResMes.SendUserId = sendId
	ontToOneResMes.Content = content
	ontToOneResData, sendErr := json.Marshal(ontToOneResMes)
	if sendErr != nil {
		sendErr = errors.New("sendMesByUserId json.Marshal(ontToOneResMes) err")
		return
	}

	mes.Data = string(ontToOneResData)
	resMes, sendErr := json.Marshal(mes)
	if sendErr != nil {
		sendErr = errors.New("sendMesByUserId json.Marshal(mes) err")
		return
	}

	sendErr = sendMes(resMes, receiveUserProcess.Conn)
	if sendErr != nil {
		sendErr = errors.New("sendMes WritePkg err")
		return
	}
	return
}

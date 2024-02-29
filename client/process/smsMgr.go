package process

import (
	"chatRoom/common/message"
	"encoding/json"
	"fmt"
)

// 管理接收的消息
type SmsManager struct {
}

var smsManager *SmsManager

func init() {
	smsManager = &SmsManager{}
}

func (this *SmsManager) showGroupMes(mesRes *message.Message) {
	var smsResMes message.SmsMes
	err := json.Unmarshal([]byte(mesRes.Data), &smsResMes)
	if err != nil {
		fmt.Println("showGroupMes json.Unmarshal err")
		return
	}
	//如果当前用户就是发送者，不做任何处理
	if currentUser.UserId == smsResMes.UserId {
		return
	}
	fmt.Println("-----您接收到一个群聊消息-----")
	mesInfo := fmt.Sprintf("用户%d : %s", smsResMes.UserId, smsResMes.Content)
	fmt.Println(mesInfo)
	fmt.Println()
}

func (this *SmsManager) showOneToOneMes(mesRead *message.Message) {
	var oneToOneResMes message.OneToOneResMes
	err := json.Unmarshal([]byte(mesRead.Data), &oneToOneResMes)
	if err != nil {
		fmt.Println("listenServerMes oneToOneResMes json.Unmarshal")
		return
	}

	fmt.Println("-----您接收到一个私聊消息-----")
	mesInfo := fmt.Sprintf("用户%d : %s", oneToOneResMes.SendUserId, oneToOneResMes.Content)
	fmt.Println(mesInfo)
	fmt.Println()
}

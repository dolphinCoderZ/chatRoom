package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	LogOutMesType           = "LogOut"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
	OneToOneMesType         = "OneToOneMes"
	OneToOneResMesType      = "OneToOneResMes"
)

const (
	UserOnline = iota
	UserOffline
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

type LoginResMes struct {
	Code          int    `json:"code"`
	OnlineUsersId []int  `json:"onlineUsersId"`
	Error         string `json:"error"`
}

type RegisterMes struct {
	//UserId   int    `json:"userId"`
	//UserPwd  string `json:"userPwd"`
	//UserName string `json:"userName"`
	User User `json:"user"`
}

type RegisterResMes struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type LogOutMes struct {
	UserId int `json:"userId"`
}

type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

type SmsMes struct {
	Content string `json:"content"`
	//匿名结构体，继承User
	User
}

type OneToOneMes struct {
	Content       string `json:"content"`
	ReceiveUserId int    `json:"receiveUserId"`
	User
}

type OneToOneResMes struct {
	SendUserId int    `json:"sendUserId"`
	Content    string `json:"content"`
}

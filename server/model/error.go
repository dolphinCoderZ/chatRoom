package model

import "errors"

var (
	ErrorUserNotExists = errors.New("用户不存在")
	ErrorExists        = errors.New("用户已经存在")
	ErrorUserPwd       = errors.New("密码不正确")
	ErrorAlreadyLogin  = errors.New("用户已经登录")
)

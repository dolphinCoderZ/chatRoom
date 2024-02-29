package model

import (
	"chatRoom/common/message"
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
)

// 数据库访问对象
type UserDAO struct {
	pool *redis.Pool
}

var GlobalUserDAO *UserDAO

func NewUserDAO(pool *redis.Pool) (userDAO *UserDAO) {
	userDAO = &UserDAO{
		pool: pool,
	}
	return userDAO
}

func getUserById(conn redis.Conn, id int) (user *User, err error) {
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			err = ErrorUserNotExists
		}
		return
	}

	//指针类型的返回值参数初始化为nil，所以必须进行初始化
	user = &User{}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		return
	}
	return user, err
}

func (this *UserDAO) LoginCheck(userId int, userPwd string) (user *User, err error) {
	conn := this.pool.Get()
	defer conn.Close()

	user, err = getUserById(conn, userId)
	if err != nil {
		return
	}

	if user.UserPwd != userPwd {
		err = ErrorUserPwd
		return
	}
	return
}

func (this *UserDAO) Register(user *message.User) (err error) {
	conn := this.pool.Get()
	defer conn.Close()

	_, err = getUserById(conn, user.UserId)
	//用户已经存在
	if err == nil {
		err = ErrorExists
		return
	}

	userInfo, err := json.Marshal(user)
	if err != nil {
		return
	}
	_, err = conn.Do("HSet", "users", user.UserId, string(userInfo))
	if err != nil {
		err = errors.New("保存用户信息失败")
		return
	}
	return
}

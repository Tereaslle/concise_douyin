package model

import (
	"errors"
	"sync"
)

// UserLogin 用户登录表，和UserInfo属于一对一关系
type UserLogin struct {
	ID         int64 `gorm:"primary_key"`
	UserInfoID int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"size:200;notnull"`
}

type UserLoginDAO struct {
}

var (
	userLoginDao  *UserLoginDAO
	userLoginOnce sync.Once
)

func NewUserLoginDao() *UserLoginDAO {
	userLoginOnce.Do(func() {
		userLoginDao = new(UserLoginDAO)
	})
	return userLoginDao
}

func (u *UserLoginDAO) QueryUserLogin(username, password string, login *UserLogin) error {
	if login == nil {
		return errors.New("结构体指针为空")
	}
	database.Where("username=? and password=?", username, password).First(login)
	if login.ID == 0 {
		return errors.New("用户不存在，账号或密码出错")
	}
	return nil
}

func (u *UserLoginDAO) IsUserExistByUsername(username string) bool {
	var userLogin UserLogin
	database.Where("username=?", username).First(&userLogin)
	if userLogin.ID == 0 {
		return false
	}
	return true
}

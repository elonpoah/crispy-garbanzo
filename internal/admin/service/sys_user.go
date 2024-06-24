package service

import (
	"crispy-garbanzo/global"
	system "crispy-garbanzo/internal/admin/models"
	"crispy-garbanzo/utils"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type UserService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {
	if nil == global.FPG_DB {
		return nil, fmt.Errorf("db not init")
	}
	var user system.SysUser
	err = global.FPG_DB.Where("username = ?", u.Username).First(&user).Error
	if err == nil {
		ok, err := utils.VerifyPassword(u.Password, user.Password)
		if nil != err {
			return nil, errors.New("服务器内部错误")
		}
		if !ok {
			global.FPG_LOG.Error("verify password error", zap.Error(err))
			return nil, errors.New("密码错误")
		}
	}
	return &user, err
}

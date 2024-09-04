package service

import (
	"crispy-garbanzo/global"
	system "crispy-garbanzo/internal/app/models"
	"crispy-garbanzo/internal/app/models/request"
	"crispy-garbanzo/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
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
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
	}
	return &user, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.FPG_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	err = global.FPG_DB.Create(&u).Error
	return u, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ChangePassword
//@description: 修改用户密码
//@param: u *model.SysUser, newPassword string
//@return: userInter *model.SysUser,err error

func (userService *UserService) ChangePassword(uid int, Password string, newPassword string) (err error) {
	var user system.SysUser
	err = global.FPG_DB.Where("id = ?", uid).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(Password, user.Password); !ok {
			return errors.New("原密码错误")
		}
	}
	user.Password = utils.BcryptHash(newPassword)
	err = global.FPG_DB.Save(&user).Error
	return err

}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserInfo
//@description: 用户信息
//@param: u *model.SysUser
//@return: userInter *model.SysUser,err error

func (userService *UserService) GetUserInfo(uid int) (userInfo *system.SysUser, err error) {
	var user system.SysUser
	err = global.FPG_DB.Where("id = ?", uid).First(&user).Error
	return &user, err

}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserDepositList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserDepositList(info request.UserDepositRecordReq, uid int) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.FPG_DB.Model(&system.Deposit{}).Where("uid = ?", uid)
	var dataList []system.Deposit
	// if info.Username != "" {
	// 	db = db.Where("username = ?", info.Username)
	// }
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&dataList).Error
	return dataList, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserWithdrawList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserWithdrawList(info request.UserWithdrawRecordReq, uid int) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.FPG_DB.Model(&system.Withdrawal{}).Where("uid = ?", uid)
	var dataList []system.Withdrawal
	// if info.Username != "" {
	// 	db = db.Where("username = ?", info.Username)
	// }
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&dataList).Error
	return dataList, total, err
}

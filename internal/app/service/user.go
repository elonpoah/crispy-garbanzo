package service

import (
	commonReq "crispy-garbanzo/common/request"
	"crispy-garbanzo/common/response"
	"crispy-garbanzo/global"
	system "crispy-garbanzo/internal/app/models"
	"crispy-garbanzo/internal/app/models/request"
	"crispy-garbanzo/utils"
	"errors"

	"gorm.io/gorm"
)

type UserService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.SysUser) (user *system.SysUser, errCode int) {
	if nil == global.FPG_DB {
		return nil, response.InternalServerError
	}
	err := global.FPG_DB.Where("username = ?", u.Username).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, response.PasswordError
		}
	}
	if err != nil {
		return nil, response.UserNotFound
	}
	return user, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

func (userService *UserService) Register(u system.SysUser) (user *system.SysUser, errCode int) {
	if !errors.Is(global.FPG_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return user, response.UserNameAlready
	}
	// 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	err := global.FPG_DB.Create(&u).Error
	if err != nil {
		return nil, response.InternalServerError
	}
	return &u, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ChangePassword
//@description: 修改用户密码
//@param: u *model.SysUser, newPassword string
//@return: userInter *model.SysUser,err error

func (userService *UserService) ChangePassword(uid int, Password string, newPassword string) (errCode int) {
	var user system.SysUser
	err := global.FPG_DB.Where("id = ?", uid).First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(Password, user.Password); !ok {
			return response.OldPasswordError
		}
	}
	if err != nil {
		return response.InternalServerError
	}
	user.Password = utils.BcryptHash(newPassword)
	err = global.FPG_DB.Save(&user).Error
	if err != nil {
		return response.InternalServerError
	}
	return response.SUCCESS

}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserInfo
//@description: 用户信息
//@param: u *model.SysUser
//@return: userInter *model.SysUser,err error

func (userService *UserService) GetUserInfo(uid int) (userInfo *system.SysUser, errCode int) {
	err := global.FPG_DB.Where("id = ?", uid).First(&userInfo).Error
	if err != nil {
		return userInfo, response.InternalServerError
	}
	return userInfo, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserDepositList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserDepositList(info request.UserDepositRecordReq, uid int) (list *[]system.Deposit, total int64, errCode int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.FPG_DB.Model(&system.Deposit{}).Where("uid = ?", uid)
	// if info.Username != "" {
	// 	db = db.Where("username = ?", info.Username)
	// }
	err := db.Count(&total).Error
	if err != nil {
		return list, total, response.InternalServerError
	}
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return list, total, response.InternalServerError
	}
	return list, total, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Deposit
//@description: 分页获取数据
//@param: info request.UserDepositRecordReq
//@return: address string, err error

func (userService *UserService) Deposit(req request.UserDepositReq) (address string, errCode int) {
	var record system.WalletAddress
	err := global.FPG_DB.Where("uid = ? AND type = ?", req.Uid, req.Type).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = global.FPG_DB.Where("enable = ? AND status = ? AND type = ?", 1, 0, req.Type).First(&record).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return address, response.ChannelCrowded
		} else if err != nil {
			return address, response.InternalServerError
		}
		record.Status = 1
		record.Uid = req.Uid
		err = global.FPG_DB.Save(&record).Error
		if err != nil {
			return address, response.InternalServerError
		}
		return record.Address, response.SUCCESS
	} else {
		return record.Address, response.SUCCESS
	}
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Deposit
//@description: 分页获取数据
//@param: info request.UserWithdrawReq
//@return: address string, err error

func (userService *UserService) Withdraw(req request.UserWithdrawReq) (errCode int) {
	withdraw := system.Withdrawal{
		Uid:       req.Uid,
		Username:  req.Username,
		Type:      req.Type,
		Amount:    req.Amount,
		ToAddress: req.Address,
	}
	err := global.FPG_DB.Create(&withdraw).Error
	if err != nil {
		return response.InternalServerError
	}
	return response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserWithdrawList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserWithdrawList(info request.UserWithdrawRecordReq, uid int) (list *[]system.Withdrawal, total int64, errCode int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.FPG_DB.Model(&system.Withdrawal{}).Where("uid = ?", uid)
	// if info.Username != "" {
	// 	db = db.Where("username = ?", info.Username)
	// }
	err := db.Count(&total).Error
	if err != nil {
		return list, total, response.InternalServerError
	}
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return list, total, response.InternalServerError
	}
	return list, total, response.SUCCESS
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserWithdrawList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserFreeSpinList(info commonReq.PageInfo, uid int) (list *[]system.InviteDuty, total int64, errCode int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.FPG_DB.Model(&system.InviteDuty{}).Where("uid = ?", uid)
	// if info.Username != "" {
	// 	db = db.Where("username = ?", info.Username)
	// }
	err := db.Count(&total).Error
	if err != nil {
		return list, total, response.InternalServerError
	}
	err = db.Limit(limit).Offset(offset).Find(&list).Error
	if err != nil {
		return list, total, response.InternalServerError
	}
	return list, total, response.SUCCESS
}

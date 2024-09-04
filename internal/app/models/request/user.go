package request

import (
	"crispy-garbanzo/common/request"
	// system "crispy-garbanzo/internal/app/models"
)

type Register struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	NickName string `json:"nickName"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

// User login structure
type Login struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
}

// Modify password structure
type ChangePasswordReq struct {
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}

type UserDepositRecordReq struct {
	request.PageInfo
	// system.Deposit
}

type UserWithdrawRecordReq struct {
	request.PageInfo
	// system.Deposit
}

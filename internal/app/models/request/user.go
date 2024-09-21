package request

import (
	"crispy-garbanzo/common/request"
	// system "crispy-garbanzo/internal/app/models"
)

type Register struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Pid      uint   `json:"inviteCode"`
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

type UserDepositReq struct {
	Uid      int     `json:"uid"`
	Username string  `json:"userName"`
	Type     int     `json:"type" binding:"required"` //2 trc20 1 erc20
	Amount   float64 `json:"amount" binding:"required"`
}

type UserWithdrawReq struct {
	Uid      int     `json:"uid"`
	Username string  `json:"userName"`
	Type     int     `json:"type" binding:"required"` //2 trc20 1 erc20
	Amount   float64 `json:"amount" binding:"required"`
	Address  string  `json:"address" binding:"required"`
}

type UserWithdrawRecordReq struct {
	request.PageInfo
	// system.Deposit
}

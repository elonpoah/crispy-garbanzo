package system

import (
	"crispy-garbanzo/global"
)

type SysUser struct {
	global.Model
	Username      string  `json:"userName" gorm:"index;comment:用户登录名"`             // 用户登录名
	Password      string  `json:"-" gorm:"type:varchar(125);comment:用户登录密码"`       // 用户登录密码
	NickName      string  `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`       // 用户昵称
	Phone         string  `json:"phone" gorm:"comment:用户手机号"`                      // 用户手机号
	Email         string  `json:"email" gorm:"comment:用户邮箱"`                       // 用户邮箱
	Enable        int     `json:"enable" gorm:"default:1;comment:用户是否被冻结 1正常 2冻结"` //用户是否被冻结 1正常 2冻结
	Balance       float64 `json:"balance" gorm:"type:decimal(12,2);"`
	FreezeBalance float64 `json:"freezeBalance" gorm:"default:0;type:decimal(12,2);"`
}

func (*SysUser) TableName() string {
	return "app_user"
}

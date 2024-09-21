package system

import (
	"crispy-garbanzo/global"
)

type InviteDuty struct {
	global.Model
	Uid    int `json:"-"`
	Type   int `json:"type" gorm:"type:tinyint;comment:1:daily 2:weekly 3:monthly"`
	Amount int `json:"amount"`
}

func (InviteDuty) TableName() string {
	return "app_invite_duty"
}

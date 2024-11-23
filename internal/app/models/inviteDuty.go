package system

import (
	"crispy-garbanzo/global"
)

type InviteDuty struct {
	global.Model
	Uid    int     `json:"-"`
	Type   int     `json:"type" gorm:"type:tinyint;comment:1:daily 2:weekly 3:monthly 4:红包"`
	Amount float64 `json:"amount" gorm:"type:decimal(12,2);"`
}

func (InviteDuty) TableName() string {
	return "app_invite_duty"
}

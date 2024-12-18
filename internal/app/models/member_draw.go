package system

import (
	"crispy-garbanzo/global"
	"time"
)

type MemberDraw struct {
	global.Model
	Uid          int       `json:"-"`
	Username     string    `json:"-"`
	DrawId       string    `json:"drawId"`
	BonusType    int       `json:"bonusType" gorm:"type:tinyint;comment:抽奖类型 1 随机 2平均"`
	Bonus        float64   `json:"bonus" gorm:"type:decimal(12,2);comment:金额"`
	Distribute   float64   `json:"distribute" gorm:"type:decimal(12,2);default:0;comment:派奖金额"`
	Count        int       `json:"count" gorm:"type:tinyint unsigned;default:1;comment:个数"`
	Participants uint      `json:"participants" gorm:"default:0;comment:参与人数"`
	Status       uint      `json:"status" gorm:"default:1;comment:状态 1:进行中, 2:已结束, 3:已过期"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"comment:过期时间"`
}

func (MemberDraw) TableName() string {
	return "app_member_draw"
}

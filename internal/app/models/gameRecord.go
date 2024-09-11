package system

import (
	"crispy-garbanzo/global"
)

type GameRecord struct {
	global.Model
	Uid                uint   `json:"-" gorm:"comment:用户ID"`
	Username           string `json:"-" gorm:"comment:用户名字"`
	SessionID          uint   `json:"sessionId" gorm:"comment:场次ID"`
	ActivytyName       string `json:"activityName" gorm:"comment:活动名称"`
	ActivytyBonus      uint   `json:"activityBonus" gorm:"comment:活动奖金"`
	ActivytySpend      uint   `json:"activitySpend" gorm:"comment:活动入场券金额"`
	ActivytyLimitCount uint   `json:"activityLimitCount" gorm:"comment:活动每场参与人数"`
	OpenTime           int64  `json:"openTime" gorm:"comment:场次时间"`
	Uids               uint   `json:"uids" gorm:"default:0;comment:参与人数"`
	Status             uint   `json:"status" gorm:"default:1;comment:状态 1:待开奖, 2:已中奖, 3:未中奖, 4:作废"`
}

func (GameRecord) TableName() string {
	return "app_game_record"
}

package system

import (
	"crispy-garbanzo/global"
)

type ActivitySession struct {
	global.Model
	ActivityID         uint   `json:"activityId" gorm:"comment:活动ID"`
	ActivytyName       string `json:"activityName" gorm:"comment:活动名称"`
	ActivytyBonus      uint   `json:"activityBonus" gorm:"comment:活动奖金"`
	ActivytySpend      uint   `json:"activitySpend" gorm:"comment:活动入场券金额"`
	ActivytyLimitCount uint   `json:"activityLimitCount" gorm:"comment:活动每场参与人数"`
	OpenTime           int64  `json:"openTime" gorm:"comment:场次时间"`
	Uids               uint   `json:"uids" gorm:"default:0;comment:参与人数"`
	Status             int8   `json:"status" gorm:"default:0;comment:状态 -1:作废, 0:待开奖, 1:已开奖"`
}

func (ActivitySession) TableName() string {
	return "app_activity_session"
}

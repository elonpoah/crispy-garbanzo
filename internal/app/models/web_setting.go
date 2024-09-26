package system

import (
	"encoding/json"

	"crispy-garbanzo/global"
)

const (
	AppWithdrawSetting = "APP_WITHDRAW_SETTING"
)

type WebStting struct {
	global.Model
	KeyName string          `json:"key"`
	Info    json.RawMessage `json:"value" gorm:"type:json"`
}

func (WebStting) TableName() string {
	return "app_setting"
}

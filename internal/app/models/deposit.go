package system

import (
	"crispy-garbanzo/global"
)

type Deposit struct {
	global.Model
	Uid         string  `json:"uid"`
	Username    string  `json:"userName"`
	Type        int     `gorm:"type:tinyint"` //0 trc20 1 erc20
	Amount      float64 `gorm:"type:decimal(12,2)"`
	FromAddress string  `json:"fromAddress"`
	ToAddress   string  `json:"toAddress"`
	TxHash      string  `gorm:"unique"`
	Status      int     `gorm:"type:tinyint unsigned;default:0"` // 0 待确认 1 成功 2 失败
}

func (Deposit) TableName() string {
	return "app_deposit"
}

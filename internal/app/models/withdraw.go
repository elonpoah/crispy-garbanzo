package system

import (
	"crispy-garbanzo/global"
)

type Withdrawal struct {
	global.Model
	Uid         int     `json:"-"`
	Username    string  `json:"-"`
	Type        int     `json:"type" gorm:"type:tinyint"` //2 trc20 1 erc20
	Amount      float64 `json:"amount" gorm:"type:decimal(12,2)"`
	FromAddress string  `json:"fromAddress"`
	ToAddress   string  `json:"toAddress"`
	TxHash      string  `json:"txHash" gorm:"unique"`
	Status      string  `json:"status" gorm:"type:tinyint unsigned;default:0"` // 0 待确认 1 成功 2 失败
}

func (Withdrawal) TableName() string {
	return "app_withdrawal"
}

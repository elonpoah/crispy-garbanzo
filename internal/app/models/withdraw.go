package system

import (
	"crispy-garbanzo/global"
)

type Withdrawal struct {
	global.Model
	Uid         int     `json:"-"`
	Username    string  `json:"-"`
	Type        int     `json:"type" gorm:"type:tinyint;comment:2 trc20 1 erc20"` //2 trc20 1 erc20
	Amount      float64 `json:"amount" gorm:"type:decimal(12,2)"`
	FromAddress string  `json:"-"`
	ToAddress   string  `json:"toAddress"`
	TxHash      string  `json:"-"`
	Status      int     `json:"status" gorm:"type:tinyint unsigned;default:0;comment:0 待确认 1 成功 2 失败"` // 0 待确认 1 成功 2 失败
	Remark      string  `json:"-" gorm:"comment:备注"`
}

func (Withdrawal) TableName() string {
	return "app_withdrawal"
}

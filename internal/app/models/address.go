package system

import (
	"time"
)

type WalletAddress struct {
	ID        uint      `gorm:"primarykey;autoIncrement;comment:主键编码" json:"ID"` // 主键ID
	Uid       int       `json:"uid"`
	CreatedAt time.Time `json:"createdAt"`                // 创建时间
	UpdatedAt time.Time `json:"-"`                        // 更新时间
	Type      int       `json:"type" gorm:"type:tinyint"` //2 trc20 1 erc20
	Enable    int       `json:"enable"`                   //1正常 2冻结
	Address   string    `json:"address"`
	Status    int       `json:"status" gorm:"type:tinyint unsigned;default:0"` // 0 空闲 1 占用中
}

func (WalletAddress) TableName() string {
	return "app_wallet_address"
}

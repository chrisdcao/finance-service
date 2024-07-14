package models

import (
	"time"
)

type Wallet struct {
	Id           uint `gorm:"primaryKey"`
	UserId       int
	Balance      float64
	WalletType   string
	CreatedOn    time.Time `gorm:"autoCreateTime"`
	UpdatedOn    time.Time `gorm:"autoUpdateTime"`
	CreatedBy    string
	UpdatedBy    string
	Transactions []Transaction `gorm:"foreignKey:WalletId"`
}

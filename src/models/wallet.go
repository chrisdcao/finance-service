package models

import (
	"time"
)

type Wallet struct {
	ID           uint          `gorm:"primaryKey"`
	UserId       int           `json:"user_id"`
	Balance      float64       `json:"balance"`
	WalletType   string        `json:"wallet_type"`
	CreatedOn    time.Time     `json:"created_on"`
	UpdatedOn    time.Time     `json:"updated_on"`
	CreatedBy    string        `json:"created_by"`
	UpdatedBy    string        `json:"updated_by"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:WalletID"`
}

package models

import (
	"time"
)

type Transaction struct {
	// autogen fields
	Id uint `gorm:"primaryKey"`
	// content fields
	WalletId        uint
	Amount          float64
	TransactionType string
	Content         string
	// audit fields
	CreatedOn time.Time `gorm:"autoCreateTime"`
	UpdatedOn time.Time `gorm:"autoUpdateTime"`
	CreatedBy string
	UpdatedBy string
}

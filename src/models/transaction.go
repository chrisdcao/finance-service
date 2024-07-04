package models

import (
	"time"
)

type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	WalletID        uint      `gorm:"wallet_id"`
	Amount          float64   `gorm:"amount"`
	TransactionType string    `gorm:"transaction_type"` // credit or debit
	Content         string    `gorm:"content"`
	CreatedOn       time.Time `gorm:"created_on"`
	UpdatedOn       time.Time `gorm:"updated_on"`
	CreatedBy       string    `gorm:"created_by"`
	UpdatedBy       string    `gorm:"updated_by"`
}

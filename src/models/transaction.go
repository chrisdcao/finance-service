package models

import (
	"time"
)

type Transaction struct {
	ID        uint      `gorm:"primaryKey"`
	WalletID  uint      `json:"wallet_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"` // credit or debit
	Content   string    `json:"content"`
	UUID      string    `json:"uuid"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}

package dto

import (
	"time"
)

type TransactionDto struct {
	ID        uint      `json:"id"`
	WalletID  uint      `json:"wallet_id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"` // credit or debit
	Content   string    `json:"content"`
	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`
	CreatedBy string    `json:"created_by"`
	UpdatedBy string    `json:"updated_by"`
}

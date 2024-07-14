package dto

import (
	"time"
)

type TransactionDto struct {
	//WalletId        uint      `json:"wallet_id"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"type"` // credit or debit
	Content         string    `json:"content"`
	CreatedOn       time.Time `json:"created_on"`
	UpdatedOn       time.Time `json:"updated_on"`
	CreatedBy       time.Time `json:"created_by"`
	UpdatedBy       time.Time `json:"updated_by"`
}

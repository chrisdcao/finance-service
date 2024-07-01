package handler

import (
	transactionDtos "finance-service/services/wallet/balance/dto"
	"gorm.io/gorm"
)

type BalanceHandler interface {
	UpdateBalance(tx *gorm.DB, input transactionDtos.UpdateBalanceInput) error
}

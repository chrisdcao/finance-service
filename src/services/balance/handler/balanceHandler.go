package handler

import (
	"context"
	transactionDtos "finance-service/services/balance/dto"
	"gorm.io/gorm"
)

type BalanceHandler interface {
	UpdateBalance(ctx context.Context, tx *gorm.DB, input transactionDtos.UpdateBalanceInput) error
}

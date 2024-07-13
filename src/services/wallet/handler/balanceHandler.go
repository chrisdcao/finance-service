package handler

import (
	"context"
	"finance-service/models"
	transactionDtos "finance-service/services/wallet/dto"
	"gorm.io/gorm"
)

type BalanceHandler interface {
	UpdateBalance(ctx context.Context, tx *gorm.DB, input transactionDtos.UpdateBalanceInput) (*models.Wallet, error)
}

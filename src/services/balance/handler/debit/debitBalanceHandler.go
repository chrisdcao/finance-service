package debit

import (
	"context"
	"finance-service/repositories"
	transactionDtos "finance-service/services/balance/dto"
	"finance-service/services/balance/handler"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DebitBalanceHandler struct {
	Repo *repositories.WalletRepository
}

func NewDebitTransaction(repo *repositories.WalletRepository) handler.BalanceHandler {
	return &DebitBalanceHandler{Repo: repo}
}

func (this *DebitBalanceHandler) UpdateBalance(ctx context.Context, tx *gorm.DB, input transactionDtos.UpdateBalanceInput) error {
	var debitAmount = -input.DiffAmount
	err := this.Repo.UpdateBalance(this.Repo.DB, nil)
	if err != nil {
		return errors.Wrap(err, "failed to update wallet balance: "+input.ToString())
	}
	return nil
}

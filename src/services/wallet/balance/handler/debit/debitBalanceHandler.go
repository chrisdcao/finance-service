package debit

import (
	"context"
	"finance-service/repositories"
	transactionDtos "finance-service/services/wallet/balance/dto"
	"finance-service/services/wallet/balance/handler"
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
	var debitAmount = -input.Amount
	err := this.Repo.UpdateBalance(this.Repo.DB, input.WalletId, debitAmount)
	if err != nil {
		return errors.Wrap(err, "failed to update wallet balance: "+input.ToString())
	}
	return nil
}

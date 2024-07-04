package credit

import (
	"context"
	"finance-service/repositories"
	transactionDtos "finance-service/services/wallet/balance/dto"
	"finance-service/services/wallet/balance/handler"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CreditBalanceHandler struct {
	Repo *repositories.WalletRepository
}

func NewCreditTransaction(repo *repositories.WalletRepository) handler.BalanceHandler {
	return &CreditBalanceHandler{Repo: repo}
}

func (this *CreditBalanceHandler) UpdateBalance(ctx context.Context, tx *gorm.DB, input transactionDtos.UpdateBalanceInput) error {
	var creditAmount = input.Amount
	err := this.Repo.UpdateBalance(tx, input.WalletId, creditAmount)
	if err != nil {
		return errors.Wrap(err, "Failed to update wallet balance"+input.ToString())
	}
	return nil
}

package credit

import (
	"context"
	"finance-service/repositories"
	transactionDtos "finance-service/services/balance/dto"
	"finance-service/services/balance/handler"
	walletservices "finance-service/services/wallet"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CreditBalanceHandler struct {
	Repo         *repositories.WalletRepository
	WalletMapper *walletservices.WalletMapper
}

func NewCreditTransaction(repo *repositories.WalletRepository, walletMapper *walletservices.WalletMapper) handler.BalanceHandler {
	return &CreditBalanceHandler{Repo: repo, WalletMapper: walletMapper}
}

func (this *CreditBalanceHandler) UpdateBalance(ctx context.Context, tx *gorm.DB, input transactionDtos.UpdateBalanceInput) error {
	wallet, err := this.WalletMapper.FromUpdateBalanceInputToWalletModel(input)
	err = this.Repo.UpdateBalance(tx, wallet)
	if err != nil {
		return errors.Wrap(err, "Failed to update wallet balance"+input.ToString())
	}
	return nil
}

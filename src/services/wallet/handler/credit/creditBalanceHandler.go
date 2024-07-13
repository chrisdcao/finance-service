package credit

import (
	"context"
	"finance-service/models"
	"finance-service/repositories"
	transactionDtos "finance-service/services/wallet/dto"
	"finance-service/services/wallet/handler"
	walletservices "finance-service/services/wallet/mapper"
	"fmt"
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

func (this *CreditBalanceHandler) UpdateBalance(ctx context.Context, tx *gorm.DB, input transactionDtos.UpdateBalanceInput) (*models.Wallet, error) {
	wallet, err := this.Repo.GetByUserIDAndWalletType(tx, input.UserId, input.WalletType.String())
	if err != nil {
		errMsg := fmt.Sprintf("Error retrieving walletId %d", wallet.ID)
		return nil, errors.Wrap(err, errMsg)
	}

	wallet.Balance = wallet.Balance + input.DiffAmount

	updatedWallet, err := this.Repo.UpdateBalance(tx, *wallet)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to update wallet balance"+input.ToString())
	}

	return updatedWallet, nil
}

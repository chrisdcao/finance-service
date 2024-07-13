package debit

import (
	"context"
	"finance-service/models"
	"finance-service/repositories"
	transactionDtos "finance-service/services/balance/dto"
	"finance-service/services/balance/handler"
	walletservices "finance-service/services/wallet/mapper"
	"finance-service/services/wallet/validator"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DebitBalanceHandler struct {
	Repo            *repositories.WalletRepository
	WalletValidator *validator.DefaultWalletValidator
	WalletMapper    *walletservices.WalletMapper
}

func NewDebitTransaction(repo *repositories.WalletRepository, mapper *walletservices.WalletMapper, validator *validator.DefaultWalletValidator) handler.BalanceHandler {
	return &DebitBalanceHandler{Repo: repo, WalletMapper: mapper, WalletValidator: validator}
}

func (this *DebitBalanceHandler) UpdateBalance(ctx context.Context, tx *gorm.DB, input transactionDtos.UpdateBalanceInput) (*models.Wallet, error) {
	wallet, err := this.Repo.GetByUserIDAndWalletType(tx, input.UserId, input.WalletType)
	if err != nil {
		errMsg := fmt.Sprintf("Error retrieving walletId %d", wallet.ID)
		return nil, errors.Wrap(err, errMsg)
	}

	isSufficientFund := this.WalletValidator.ValidateDebitWallet(*wallet, input.DiffAmount)
	if !isSufficientFund {
		errMsg := fmt.Sprintf("Current walletId %d has insufficient fund for amount: %f", wallet.ID, input.DiffAmount)
		return nil, errors.Wrap(err, errMsg)
	}

	wallet.Balance = wallet.Balance - input.DiffAmount
	updatedWallet, err := this.Repo.UpdateBalance(tx, *wallet)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to update wallet balance"+input.ToString())
	}

	return updatedWallet, nil
}

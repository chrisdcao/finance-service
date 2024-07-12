package validator

import (
	"context"
	"finance-service/models"
	walletservices "finance-service/services/wallet/read"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DefaultWalletValidator struct {
	WalletReadService *walletservices.DefaultWalletReadService
}

func NewWalletValidator(defaultWalletReadService *walletservices.DefaultWalletReadService) *DefaultWalletValidator {
	return &DefaultWalletValidator{
		WalletReadService: defaultWalletReadService,
	}
}

func (this *DefaultWalletValidator) ValidateWallets(ctx context.Context, tx *gorm.DB, toExternalWalletId string, fromExternalWalletId string, amount float64) error {
	_, err := this.WalletReadService.GetFromExternalId(ctx, toExternalWalletId, tx)
	if err != nil {
		return err
	}

	fromWallet, err := this.WalletReadService.GetFromExternalId(ctx, fromExternalWalletId, tx)
	if err != nil {
		return err
	}

	return this.validateWalletRemainingBalance(ctx, fromWallet, amount)
}

func (this *DefaultWalletValidator) ValidateTransferAmount(ctx context.Context, amount float64) error {
	// TODO: Configure amount to be a constant in a centralized file
	if amount < 2000000 {
		return errors.New("transfer amount must be greater than 2.000.000")
	}

	return nil
}

func (this *DefaultWalletValidator) validateWalletRemainingBalance(ctx context.Context, fromWallet *models.Wallet, amount float64) error {
	if fromWallet.Balance < amount {
		return errors.New("insufficient wallet balance: " + fmt.Sprintf("%.2f", fromWallet.Balance))
	}

	return nil
}

package validator

import (
	"context"
	"finance-service/models"
	"github.com/pkg/errors"
)

type DefaultWalletValidator struct{}

func NewWalletValidator() *DefaultWalletValidator {
	return &DefaultWalletValidator{}
}

func (this *DefaultWalletValidator) ValidateTransferAmount(ctx context.Context, amount float64) error {
	// TODO: Configure amount to be a constant in a centralized file
	if amount < 2000000 {
		return errors.New("transfer amount must be greater than 2.000.000")
	}

	return nil
}

func (this *DefaultWalletValidator) ValidateDebitWallet(wallet models.Wallet, amountToDebit float64) bool {
	return wallet.Balance-amountToDebit > 0.0
}

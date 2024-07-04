package dto

import (
	"finance-service/services/wallet/balance/enums"
	"fmt"
)

type UpdateBalanceInput struct {
	WalletId         uint
	WalletType       string
	Amount           float64
	BalanceOperation enums.BalanceOperation
}

// ToString method to return a string representation of UpdateBalanceInput
func (this *UpdateBalanceInput) ToString() string {
	return fmt.Sprintf(
		"UpdateBalanceInput{WalletId: %d, WalletType: %s, Amount: %.2f, BalanceOperation: %s}",
		this.WalletId, this.WalletType, this.Amount, this.BalanceOperation,
	)
}

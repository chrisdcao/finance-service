package dto

import (
	"finance-service/services/balance/enums"
	"fmt"
)

type UpdateBalanceInput struct {
	UserId           string
	WalletType       string
	DiffAmount       float64
	BalanceOperation enums.BalanceOperation
	Content          string
}

// ToString method to return a string representation of UpdateBalanceInput
func (this *UpdateBalanceInput) ToString() string {
	return fmt.Sprintf(
		"UpdateBalanceInput{UserId: %s, WalletType: %s, DiffAmount: %.2f, BalanceOperation: %s}",
		this.UserId, this.WalletType, this.DiffAmount, this.BalanceOperation,
	)
}

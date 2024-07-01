package dto

import (
	"finance-service/services/wallet/balance/enums"
)

type UpdateBalanceInput struct {
	WalletId         uint
	WalletType       string
	Amount           float64
	BalanceOperation enums.BalanceOperation
}

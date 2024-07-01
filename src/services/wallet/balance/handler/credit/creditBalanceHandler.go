package credit

import (
	"finance-service/repositories"
	transactionDtos "finance-service/services/wallet/balance/dto"
	"finance-service/services/wallet/balance/handler"
	"gorm.io/gorm"
)

type CreditBalanceHandler struct {
	Repo *repositories.WalletRepository
}

func NewCreditTransaction(repo *repositories.WalletRepository) handler.BalanceHandler {
	return &CreditBalanceHandler{Repo: repo}
}

func (this *CreditBalanceHandler) UpdateBalance(tx *gorm.DB, input transactionDtos.UpdateBalanceInput) error {
	var creditAmount = input.Amount
	return this.Repo.UpdateBalance(tx, input.WalletId, creditAmount)
}

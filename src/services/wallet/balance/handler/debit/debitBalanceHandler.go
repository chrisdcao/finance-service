package debit

import (
	"finance-service/repositories"
	transactionDtos "finance-service/services/wallet/balance/dto"
	"finance-service/services/wallet/balance/handler"
	"gorm.io/gorm"
)

type DebitBalanceHandler struct {
	Repo *repositories.WalletRepository
}

func NewDebitTransaction(repo *repositories.WalletRepository) handler.BalanceHandler {
	return &DebitBalanceHandler{Repo: repo}
}

func (this *DebitBalanceHandler) UpdateBalance(tx *gorm.DB, input transactionDtos.UpdateBalanceInput) error {
	var debitAmount = -input.Amount
	return this.Repo.UpdateBalance(this.Repo.DB, input.WalletId, debitAmount)
}

package debit

import (
	"finance-service/repositories"
	transactionDtos "finance-service/services/wallet/balance/dto"
	"finance-service/services/wallet/balance/handler"
)

type DebitBalanceHandler struct {
	Repo *repositories.WalletRepository
}

func NewDebitTransaction(repo *repositories.WalletRepository) handler.BalanceHandler {
	return &DebitBalanceHandler{Repo: repo}
}

func (t *DebitBalanceHandler) UpdateBalance(input transactionDtos.UpdateBalanceInput) error {
	var debitAmount = -input.Amount
	return t.Repo.UpdateBalance(input.WalletId, debitAmount)
}

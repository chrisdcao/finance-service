package response

import (
	transactionDtos "finance-service/services/transaction/dto"
)

type WalletTransferResponse struct {
	Transactions []transactionDtos.TransactionDto
}

func NewWalletTransferResponse(transactions []transactionDtos.TransactionDto) WalletTransferResponse {
	return WalletTransferResponse{Transactions: transactions}
}

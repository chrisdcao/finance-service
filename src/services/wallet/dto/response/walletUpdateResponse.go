package response

import (
	transactionDtos "finance-service/services/transaction/dto"
)

type WalletUpdateResponse struct {
	Transaction transactionDtos.TransactionDto
}

func NewWalletUpdateResponse(transaction transactionDtos.TransactionDto) *WalletUpdateResponse {
	return &WalletUpdateResponse{
		Transaction: transaction,
	}
}

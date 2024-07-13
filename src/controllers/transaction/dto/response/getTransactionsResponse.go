package response

import transactionDtos "finance-service/services/transaction/dto"

type GetTransactionsResponse struct {
	Transactions []transactionDtos.TransactionDto
}

func NewGetTransactionsResponse(transactions []transactionDtos.TransactionDto) GetTransactionsResponse {
	return GetTransactionsResponse{Transactions: transactions}
}

package mapper

import (
	"finance-service/models"
	"finance-service/services/wallet/transaction/dto"
)

type TransactionDtoMapper struct {
}

func (this *TransactionDtoMapper) ToTransactionModel(dto dto.TransactionDto) models.Transaction {
	return models.Transaction{
		ID:       dto.ID,
		WalletID: dto.WalletID,
		Amount:   dto.Amount,
		Type:     dto.Type,
		Content:  dto.Content,
	}
}

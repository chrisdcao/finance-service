package mapper

import (
	"finance-service/models"
	"finance-service/services/wallet/transaction/dto"
)

type TransactionMapper struct {
}

func (this *TransactionMapper) FromDtoToModel(dto dto.TransactionDto) models.Transaction {
	return models.Transaction{
		ID:       dto.ID,
		WalletID: dto.WalletID,
		Amount:   dto.Amount,
		Type:     dto.Type,
		Content:  dto.Content,
	}
}

func (this *TransactionMapper) FromModelToDto(transaction models.Transaction) dto.TransactionDto {
	return dto.TransactionDto{
		ID:       transaction.ID,
		WalletID: transaction.WalletID,
		Amount:   transaction.Amount,
		Type:     transaction.Type,
		Content:  transaction.Content,
	}
}

func (this *TransactionMapper) FromModelListToDtoList(transactions []models.Transaction) []dto.TransactionDto {
	var transactionDtos []dto.TransactionDto
	for _, transaction := range transactions {
		transactionDto := this.FromModelToDto(transaction)
		transactionDtos = append(transactionDtos, transactionDto)
	}
	return transactionDtos
}

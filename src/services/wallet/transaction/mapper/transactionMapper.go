package mapper

import (
	"finance-service/models"
	walletDtos "finance-service/services/wallet/dto"
	"finance-service/services/wallet/transaction/dto"
)

type TransactionMapper struct {
}

func NewTransactionMapper() *TransactionMapper {
	return &TransactionMapper{}
}

func (this *TransactionMapper) FromDtoToModel(dto dto.TransactionDto) models.Transaction {
	return models.Transaction{
		ID:              dto.ID,
		WalletID:        dto.WalletID,
		Amount:          dto.Amount,
		TransactionType: dto.TransactionType,
		Content:         dto.Content,
	}
}

func (this *TransactionMapper) FromModelToDto(transaction models.Transaction) dto.TransactionDto {
	return dto.TransactionDto{
		ID:              transaction.ID,
		WalletID:        transaction.WalletID,
		Amount:          transaction.Amount,
		TransactionType: transaction.TransactionType,
		Content:         transaction.Content,
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

func (this *TransactionMapper) FromWalletIdAndRequesToDto(walletId uint, updateRequest walletDtos.WalletUpdateRequest) dto.TransactionDto {
	return dto.TransactionDto{
		WalletID:        walletId,
		Amount:          updateRequest.Amount,
		TransactionType: updateRequest.UpdateType,
		Content:         updateRequest.Content,
	}
}

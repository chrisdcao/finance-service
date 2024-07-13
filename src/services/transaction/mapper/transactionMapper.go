package mapper

import (
	"finance-service/models"
	"finance-service/services/transaction/dto"
	transactionDtos "finance-service/services/wallet/dto"
)

type TransactionMapper struct {
}

func NewTransactionMapper() *TransactionMapper {
	return &TransactionMapper{}
}

func (this *TransactionMapper) FromDtoToModel(dto dto.TransactionDto) models.Transaction {
	return models.Transaction{
		Amount:          dto.Amount,
		TransactionType: dto.TransactionType,
		Content:         dto.Content,
	}
}

func (this *TransactionMapper) FromModelToDto(transaction models.Transaction) *dto.TransactionDto {
	return &dto.TransactionDto{
		Amount:          transaction.Amount,
		TransactionType: transaction.TransactionType,
		Content:         transaction.Content,
	}
}

func (this *TransactionMapper) FromModelListToDtoList(transactions []models.Transaction) []dto.TransactionDto {
	var transactionDtos []dto.TransactionDto
	for _, transaction := range transactions {
		transactionDto := this.FromModelToDto(transaction)
		transactionDtos = append(transactionDtos, *transactionDto)
	}
	return transactionDtos
}

func (this *TransactionMapper) FromInputAndWalletToTx(updateInput transactionDtos.UpdateBalanceInput, wallet models.Wallet) *models.Transaction {
	return &models.Transaction{
		WalletID:        wallet.ID,
		Amount:          updateInput.DiffAmount,
		TransactionType: updateInput.BalanceOperation.String(),
		Content:         updateInput.Content,
	}
}

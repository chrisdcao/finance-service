package transaction

import (
	"finance-service/repositories"
	"finance-service/services/wallet/transaction/dto"
	"finance-service/services/wallet/transaction/mapper"
	"finance-service/utils"
	"gorm.io/gorm"
)

type TransactionWriteService struct {
	TransactionDtoMapper  *mapper.TransactionDtoMapper
	TransactionRepository *repositories.TransactionRepository
}

func NewTransactionWriteService() *TransactionWriteService {
	return &TransactionWriteService{TransactionRepository: repositories.NewTransactionRepository()}
}

func (this *TransactionWriteService) CreateTransaction(tx *gorm.DB, transactionDto dto.TransactionDto) error {
	transaction := this.TransactionDtoMapper.ToTransactionModel(transactionDto)
	if err := this.TransactionRepository.Create(tx, &transaction); err != nil {
		utils.Logger().Println("Error creating transaction:", err)
		return err
	}
	return nil
}

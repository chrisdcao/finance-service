package transaction

import (
	"finance-service/models"
	"finance-service/repositories"
	"finance-service/services/transaction/dto"
	"finance-service/services/transaction/mapper"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TransactionWriteService struct {
	TransactionMapper     *mapper.TransactionMapper
	TransactionRepository *repositories.TransactionRepository
}

func NewTransactionWriteService(transactionRepository *repositories.TransactionRepository) *TransactionWriteService {
	return &TransactionWriteService{TransactionRepository: transactionRepository}
}

func (this *TransactionWriteService) CreateTransaction(tx *gorm.DB, transaction models.Transaction) (*dto.TransactionDto, error) {
	savedTransaction, err := this.TransactionRepository.Create(tx, transaction)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transaction")
	}
	return this.TransactionMapper.FromModelToDto(*savedTransaction), nil
}

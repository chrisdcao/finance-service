package transaction

import (
	"finance-service/repositories"
	"finance-service/services/wallet/transaction/dto"
	"finance-service/services/wallet/transaction/mapper"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TransactionWriteService struct {
	TransactionDtoMapper  *mapper.TransactionMapper
	TransactionRepository *repositories.TransactionRepository
}

func NewTransactionWriteService(transactionRepository *repositories.TransactionRepository) *TransactionWriteService {
	return &TransactionWriteService{TransactionRepository: transactionRepository}
}

func (this *TransactionWriteService) CreateTransaction(tx *gorm.DB, transactionDto dto.TransactionDto) error {
	transaction := this.TransactionDtoMapper.FromDtoToModel(transactionDto)
	if err := this.TransactionRepository.Create(tx, &transaction); err != nil {
		return errors.Wrap(err, "failed to create transaction")
	}
	return nil
}

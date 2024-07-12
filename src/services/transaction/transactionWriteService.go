package transaction

import (
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

func (this *TransactionWriteService) CreateTransaction(tx *gorm.DB, transactionDto dto.TransactionDto) (*dto.TransactionDto, error) {
	transaction := this.TransactionMapper.FromDtoToModel(transactionDto)
	if err := this.TransactionRepository.Create(tx, &transaction); err != nil {
		return nil, errors.Wrap(err, "failed to create transaction")
	}
	return this.TransactionMapper.FromModelToDto(transaction), nil
}

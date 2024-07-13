package transaction

import (
	"finance-service/common/exception"
	dto2 "finance-service/controllers/transaction/dto/request"
	"finance-service/repositories"
	"finance-service/services/transaction/dto"
	"finance-service/services/transaction/mapper"
	"github.com/pkg/errors"
)

type TransactionReadService struct {
	TransactionDtoMapper  *mapper.TransactionMapper
	TransactionRepository *repositories.TransactionRepository
}

func NewTransactionReadService(repository *repositories.TransactionRepository) *TransactionReadService {
	return &TransactionReadService{TransactionRepository: repository}
}

func (this *TransactionReadService) GetTransactions(params dto2.GetTransactionsRequest) ([]dto.TransactionDto, error) {
	foundTransactions, err := this.TransactionRepository.FindTransactions(
		params.WalletType,
		params.ActionType,
		params.Amount,
		params.FromTime,
		params.ToTime,
		params.UUID,
	)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get transactions from repository!")
	}
	if len(foundTransactions) == 0 {
		return nil, errors.Wrap(exception.ErrTransactionNotFound, "No transaction found for given params!")
	}

	return this.TransactionDtoMapper.FromModelListToDtoList(foundTransactions), nil
}

package transaction

import (
	"finance-service/repositories"
	"finance-service/services/wallet/transaction/mapper"
	"finance-service/utils"
)

type TransactionReadService struct {
	TransactionDtoMapper  *mapper.TransactionDtoMapper
	TransactionRepository *repositories.TransactionRepository
}

func NewTransactionReadService() *TransactionReadService {
	return &TransactionReadService{TransactionRepository: repositories.NewTransactionRepository()}
}

func (this *TransactionReadService) GetTransactions(uuid string) ([]dto.TransactionDto, error) {
	transactions, err := this.TransactionRepository.GetAllByUuid(uuid)
	if err != nil {
		utils.Logger().Println("Error fetching transactions:", err)
		return nil, err
	}

	utils.Logger().Println("Fetched user transactions", uuid, transactions)
	return transactions, nil
}

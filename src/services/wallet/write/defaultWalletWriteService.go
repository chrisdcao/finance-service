package write

import (
	"context"
	txManagement "finance-service/config/transaction"
	"finance-service/repositories"
	balanceDtos "finance-service/services/balance/dto"
	balanceHandlerFactory "finance-service/services/balance/factory"
	mapper2 "finance-service/services/balance/mapper"
	"finance-service/services/transaction"
	"finance-service/services/transaction/dto"
	"finance-service/services/transaction/mapper"
	walletDtos "finance-service/services/wallet/dto/request"
	"finance-service/services/wallet/parser"
	"finance-service/services/wallet/validator"
	"finance-service/utils/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DefaultWalletWriteService struct {
	BalanceHandlerFactory   *balanceHandlerFactory.BalanceHandlerFactory
	WalletRepository        *repositories.WalletRepository
	TransactionWriteService *transaction.TransactionWriteService
	TransactionMapper       *mapper.TransactionMapper
	BalanceMapper           *mapper2.BalanceMapper
	WalletValidator         *validator.DefaultWalletValidator
	WalletIdParser          *parser.WalletIdParser
	Logger                  *log.CommonLogger
}

func NewWalletWriteService(
	balanceHandlerFactory *balanceHandlerFactory.BalanceHandlerFactory,
	walletRepository *repositories.WalletRepository,
	transactionWriteService *transaction.TransactionWriteService,
	transactionMapper *mapper.TransactionMapper,
	balanceMapper *mapper2.BalanceMapper,
	walletValidator *validator.DefaultWalletValidator,
	walletIdParser *parser.WalletIdParser,
) *DefaultWalletWriteService {
	return &DefaultWalletWriteService{
		BalanceHandlerFactory:   balanceHandlerFactory,
		WalletRepository:        walletRepository,
		TransactionWriteService: transactionWriteService,
		TransactionMapper:       transactionMapper,
		WalletValidator:         walletValidator,
		WalletIdParser:          walletIdParser,
		Logger:                  log.NewCommonLogger(),
		BalanceMapper:           balanceMapper,
	}
}

func (this *DefaultWalletWriteService) UpdateBalance(ctx context.Context, tx *gorm.DB, updateRequest walletDtos.WalletUpdateRequest) (*dto.TransactionDto, error) {
	var transactionDto *dto.TransactionDto

	updateInput, err := this.BalanceMapper.FromUpdateRequestToInput(updateRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to map to update input")
	}

	err = txManagement.WithTransaction(this.WalletRepository.DB, tx, func(localTx *gorm.DB) error {
		err = this.updateWalletBalance(ctx, localTx, *updateInput)
		if err != nil {
			return errors.Wrap(err, "failed to update wallet balance")
		}

		transactionDto = this.TransactionMapper.FromUpdateBalanceInputToDto(*updateInput)

		if transactionDto, err = this.TransactionWriteService.CreateTransaction(localTx, *transactionDto); err != nil {
			return errors.Wrap(err, "failed to create transaction")
		}

		return nil
	})

	return transactionDto, errors.Wrap(err, "failed to update wallet balance")
}

func (this *DefaultWalletWriteService) updateWalletBalance(ctx context.Context, tx *gorm.DB, updateInput balanceDtos.UpdateBalanceInput) error {
	balanceHandler, err := this.BalanceHandlerFactory.GetHandler(updateInput)
	if err != nil {
		return errors.Wrap(err, "failed to get balance handler")
	}

	err = balanceHandler.UpdateBalance(ctx, tx, updateInput)

	return errors.Wrap(err, "failed to update wallet balance")
}

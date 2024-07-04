package services

import (
	"context"
	txManagement "finance-service/config/transaction"
	"finance-service/repositories"
	balanceDtos "finance-service/services/wallet/balance/dto"
	operationTypes "finance-service/services/wallet/balance/enums"
	balanceHandlerFactory "finance-service/services/wallet/balance/factory"
	walletDtos "finance-service/services/wallet/dto"
	"finance-service/services/wallet/parser"
	"finance-service/services/wallet/transaction"
	"finance-service/services/wallet/transaction/mapper"
	"finance-service/services/wallet/validator"
	"finance-service/utils/log"
	"gorm.io/gorm"
)

type DefaultWalletWriteService struct {
	BalanceHandlerFactory   *balanceHandlerFactory.BalanceHandlerFactory
	WalletRepository        *repositories.WalletRepository
	TransactionWriteService *transaction.TransactionWriteService
	TransactionMapper       *mapper.TransactionMapper
	WalletValidator         *validator.DefaultWalletValidator
	WalletIdParser          *parser.WalletIdParser
	Logger                  *log.CommonLogger
}

// TODO: Should we init all of these beans outside (for singleton) then pass it into the constructor here instead?
func NewWalletWriteService(
	balanceHandlerFactory *balanceHandlerFactory.BalanceHandlerFactory,
	walletRepository *repositories.WalletRepository,
	transactionWriteService *transaction.TransactionWriteService,
	transactionMapper *mapper.TransactionMapper,
	// TODO: Add this to the container initialization
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
	}
}

func (this *DefaultWalletWriteService) UpdateBalance(ctx context.Context, tx *gorm.DB, updateRequest walletDtos.WalletUpdateRequest) (uint, error) {
	var walletId uint

	err := txManagement.WithTransaction(this.WalletRepository.DB, tx, func(localTx *gorm.DB) error {
		operationType, err := this.getOperationType(updateRequest.UpdateType)

		if err != nil {
			return err
		}

		walletId, err = this.WalletIdParser.ParseFromEncryption(ctx, updateRequest.ExternalWalletId)
		if err != nil {
			return err
		}

		if err := this.updateWalletBalance(ctx, localTx, walletId, operationType, updateRequest); err != nil {
			return err
		}

		transactionDto := this.TransactionMapper.FromWalletIdAndRequesToDto(walletId, updateRequest)

		if err := this.TransactionWriteService.CreateTransaction(localTx, transactionDto); err != nil {
			return err
		}

		return nil
	})

	return walletId, err
}

func (this *DefaultWalletWriteService) getOperationType(topupType string) (operationTypes.BalanceOperation, error) {
	return operationTypes.ParseTopupType(topupType)
}

func (this *DefaultWalletWriteService) updateWalletBalance(ctx context.Context, tx *gorm.DB, walletId uint, operationType operationTypes.BalanceOperation, updateRequest walletDtos.WalletUpdateRequest) error {
	balanceHandler, err := this.BalanceHandlerFactory.GetHandler(operationType)
	if err != nil {
		return err
	}

	return balanceHandler.UpdateBalance(ctx, tx, balanceDtos.UpdateBalanceInput{
		WalletId:         walletId,
		WalletType:       updateRequest.WalletType,
		Amount:           updateRequest.Amount,
		BalanceOperation: operationType,
	})
}

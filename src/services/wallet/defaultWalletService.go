package services

import (
	"context"
	"database/sql"
	"finance-service/config"
	txManagement "finance-service/config/transaction"
	operationTypes "finance-service/services/wallet/balance/enums"
	balanceHandlerFactory "finance-service/services/wallet/balance/factory"
	walletDtos "finance-service/services/wallet/dto"
	"finance-service/services/wallet/transaction"
	"finance-service/services/wallet/transaction/mapper"
	"finance-service/services/wallet/validator"
	"fmt"
	"gorm.io/gorm"
)

type DefaultWalletService struct {
	BalanceHandlerFactory   *balanceHandlerFactory.BalanceHandlerFactory
	TransactionWriteService *transaction.TransactionWriteService
	TransactionMapper       *mapper.TransactionMapper
	WalletValidator         *validator.DefaultWalletValidator
	WalletReadService       *DefaultWalletReadService
	WalletWriteService      *DefaultWalletWriteService
}

func NewWalletService(
	balanceHandlerFactory *balanceHandlerFactory.BalanceHandlerFactory,
	transactionWriteService *transaction.TransactionWriteService,
	transactionMapper *mapper.TransactionMapper,
	// TODO: Add this to the container initialization
	walletValidator *validator.DefaultWalletValidator,
	walletReadService *DefaultWalletReadService,
	walletWriteService *DefaultWalletWriteService,
) *DefaultWalletService {
	return &DefaultWalletService{
		BalanceHandlerFactory:   balanceHandlerFactory,
		TransactionWriteService: transactionWriteService,
		TransactionMapper:       transactionMapper,
		WalletValidator:         walletValidator,
		WalletReadService:       walletReadService,
		WalletWriteService:      walletWriteService,
	}
}

func (this *DefaultWalletService) UpdateBalance(ctx context.Context, tx *gorm.DB, updateRequest walletDtos.WalletUpdateRequest) (*walletDtos.WalletDto, error) {
	walletId, err := this.WalletWriteService.UpdateBalance(ctx, tx, updateRequest)
	if err != nil {
		return nil, err
	}

	result, err := this.WalletReadService.GetWallet(ctx, tx, walletId)
	return result, err
}

func (this *DefaultWalletService) WalletTransfer(ctx context.Context, walletTransferRequest walletDtos.WalletTransferRequest) ([]walletDtos.WalletDto, error) {
	amount, toExternalWalletId, fromExternalWalletId := walletTransferRequest.Amount, walletTransferRequest.ExternalToWalletId, walletTransferRequest.ExternalFromWalletId

	err := this.WalletValidator.ValidateTransferAmount(ctx, amount)
	if err != nil {
		return nil, err
	}

	// Begin new transaction with desired isolation level (REPEATABLE READ or SERIALIZABLE)
	var updatedVndWallet, updatedAsmWallet *walletDtos.WalletDto
	err = txManagement.WithNewTransaction(config.DB, sql.LevelRepeatableRead, func(tx *gorm.DB) error {
		err = this.WalletValidator.ValidateWallets(ctx, tx, toExternalWalletId, fromExternalWalletId, amount)
		if err != nil {
			return err
		}

		// Update `from` wallet (debit)
		updatedVndWallet, err = this.UpdateBalance(ctx, tx, walletDtos.WalletUpdateRequest{
			ExternalWalletId: fromExternalWalletId,
			UpdateType:       operationTypes.VNDWalletDebit.String(),
			Amount:           amount,
			Content:          "Chuyen tien tu VND Wallet sang ASM Wallet",
		})
		if err != nil {
			return err
		}

		// Update `to` wallet (credit)
		updatedAsmWallet, err = this.UpdateBalance(ctx, tx, walletDtos.WalletUpdateRequest{
			ExternalWalletId: toExternalWalletId,
			UpdateType:       operationTypes.ASMWalletTopup.String(),
			Amount:           amount,
			Content:          "Nhan tien tu VND Wallet",
		})
		if err != nil {
			return fmt.Errorf("failed to update to wallet: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Return both wallets
	return []walletDtos.WalletDto{*updatedVndWallet, *updatedAsmWallet}, nil
}

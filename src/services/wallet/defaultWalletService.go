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
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "failed to update wallet balance")
	}

	result, err := this.WalletReadService.GetWallet(ctx, tx, walletId)
	return result, errors.Wrap(err, "failed to get wallet")
}

func (this *DefaultWalletService) WalletTransfer(ctx context.Context, walletTransferRequest walletDtos.WalletTransferRequest) ([]walletDtos.WalletDto, error) {
	amount, toExternalWalletId, fromExternalWalletId := walletTransferRequest.Amount, walletTransferRequest.ExternalToWalletId, walletTransferRequest.ExternalFromWalletId

	err := this.WalletValidator.ValidateTransferAmount(ctx, amount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate transfer amount")
	}

	// Begin new transaction with desired isolation level (REPEATABLE READ or SERIALIZABLE)
	var updatedVndWallet, updatedAsmWallet *walletDtos.WalletDto
	err = txManagement.WithNewTransaction(config.DB, sql.LevelRepeatableRead, func(tx *gorm.DB) error {
		err = this.WalletValidator.ValidateWallets(ctx, tx, toExternalWalletId, fromExternalWalletId, amount)
		if err != nil {
			return errors.Wrap(err, "failed to validate wallets")
		}

		// Update `from` wallet (debit)
		updatedVndWallet, err = this.UpdateBalance(ctx, tx, walletDtos.WalletUpdateRequest{
			ExternalWalletId: fromExternalWalletId,
			UpdateType:       operationTypes.VNDWalletDebit.String(),
			Amount:           amount,
			Content:          "Chuyen tien tu VND Wallet sang ASM Wallet",
		})
		if err != nil {
			return errors.Wrap(err, "failed to update from wallet")
		}

		// Update `to` wallet (credit)
		updatedAsmWallet, err = this.UpdateBalance(ctx, tx, walletDtos.WalletUpdateRequest{
			ExternalWalletId: toExternalWalletId,
			UpdateType:       operationTypes.ASMWalletTopup.String(),
			Amount:           amount,
			Content:          "Nhan tien tu VND Wallet",
		})
		if err != nil {
			return errors.Wrap(err, "failed to update to wallet")
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to transfer wallet")
	}

	// Return both wallets
	return []walletDtos.WalletDto{*updatedVndWallet, *updatedAsmWallet}, nil
}

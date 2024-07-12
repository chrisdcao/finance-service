package services

import (
	"context"
	"database/sql"
	"finance-service/config"
	txManagement "finance-service/config/transaction"
	operationTypes "finance-service/services/balance/enums"
	balanceHandlerFactory "finance-service/services/balance/factory"
	"finance-service/services/transaction"
	"finance-service/services/transaction/mapper"
	walletDtos "finance-service/services/wallet/dto"
	"finance-service/services/wallet/dto/request"
	response2 "finance-service/services/wallet/dto/response"
	"finance-service/services/wallet/read"
	"finance-service/services/wallet/validator"
	"finance-service/services/wallet/write"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DefaultWalletService struct {
	BalanceHandlerFactory   *balanceHandlerFactory.BalanceHandlerFactory
	TransactionWriteService *transaction.TransactionWriteService
	TransactionMapper       *mapper.TransactionMapper
	WalletValidator         *validator.DefaultWalletValidator
	WalletReadService       *read.DefaultWalletReadService
	WalletWriteService      *write.DefaultWalletWriteService
}

func NewWalletService(
	balanceHandlerFactory *balanceHandlerFactory.BalanceHandlerFactory,
	transactionWriteService *transaction.TransactionWriteService,
	transactionMapper *mapper.TransactionMapper,
	walletValidator *validator.DefaultWalletValidator,
	walletReadService *read.DefaultWalletReadService,
	walletWriteService *write.DefaultWalletWriteService,
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

func (this *DefaultWalletService) UpdateBalance(ctx context.Context, tx *gorm.DB, updateRequest request.WalletUpdateRequest) (*response2.WalletUpdateResponse, error) {
	transactionDto, err := this.WalletWriteService.UpdateBalance(ctx, tx, updateRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update wallet balance")
	}
	return response2.NewWalletUpdateResponse(*transactionDto), nil
}

func (this *DefaultWalletService) WalletTransfer(ctx context.Context, walletTransferRequest request.WalletTransferRequest) ([]walletDtos.WalletDto, error) {
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
		_, err = this.UpdateBalance(ctx, tx, request.WalletUpdateRequest{
			UserId:     fromExternalWalletId,
			UpdateType: operationTypes.VNDWalletDebit.String(),
			Amount:     amount,
			Content:    "Chuyen tien tu VND Wallet sang ASM Wallet",
		})
		if err != nil {
			return errors.Wrap(err, "failed to update from wallet")
		}

		// Update `to` wallet (credit)
		_, err = this.UpdateBalance(ctx, tx, request.WalletUpdateRequest{
			UserId:     toExternalWalletId,
			UpdateType: operationTypes.ASMWalletTopup.String(),
			Amount:     amount,
			Content:    "Nhan tien tu VND Wallet",
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

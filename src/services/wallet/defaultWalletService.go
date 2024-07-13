package services

import (
	"context"
	"database/sql"
	"finance-service/configs"
	txManagement "finance-service/configs/transaction"
	"finance-service/controllers/wallet/dto/request"
	"finance-service/services/transaction"
	"finance-service/services/transaction/dto"
	"finance-service/services/transaction/mapper"
	walletTypes "finance-service/services/wallet/enums"
	balanceHandlerFactory "finance-service/services/wallet/factory"
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

func (this *DefaultWalletService) UpdateBalance(ctx context.Context, tx *gorm.DB, updateRequest request.WalletUpdateRequest) (*dto.TransactionDto, error) {
	transactionDto, err := this.WalletWriteService.UpdateBalance(ctx, tx, updateRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update wallet balance")
	}
	return transactionDto, nil
}

func (this *DefaultWalletService) WalletTransfer(ctx context.Context, transferRequest request.WalletTransferRequest) ([]dto.TransactionDto, error) {
	err := this.WalletValidator.ValidateTransferAmount(ctx, transferRequest.Amount)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate transfer amount")
	}

	// Begin new transaction with desired isolation level (REPEATABLE READ or SERIALIZABLE)
	var tx1, tx2 *dto.TransactionDto
	err = txManagement.WithNewTransaction(configs.DB, sql.LevelRepeatableRead, func(tx *gorm.DB) error {
		if err != nil {
			return errors.Wrap(err, "failed to validate wallets")
		}

		// Update `from` wallet (debit)
		tx1, err = this.UpdateBalance(ctx, tx, request.WalletUpdateRequest{
			UserId:     transferRequest.UserId,
			WalletType: walletTypes.VNDWallet.String(),
			UpdateType: walletTypes.Debit.String(),
			Amount:     transferRequest.Amount,
			Content:    "Chuyen tien tu VND Wallet sang ASM Wallet",
		})
		if err != nil {
			return errors.Wrap(err, "failed to update from wallet")
		}

		// Update `to` wallet (credit)
		tx2, err = this.UpdateBalance(ctx, tx, request.WalletUpdateRequest{
			UserId:     transferRequest.UserId,
			WalletType: walletTypes.ASMWallet.String(),
			UpdateType: walletTypes.Credit.String(),
			Amount:     transferRequest.Amount,
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
	return []dto.TransactionDto{*tx1, *tx2}, nil
}

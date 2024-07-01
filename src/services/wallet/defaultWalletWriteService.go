package services

import (
	"database/sql"
	"errors"
	txManagement "finance-service/config/transaction"
	"finance-service/models"
	"finance-service/repositories"
	"finance-service/services/cryptography"
	balanceDtos "finance-service/services/wallet/balance/dto"
	operationTypes "finance-service/services/wallet/balance/enums"
	balanceHandlerFactory "finance-service/services/wallet/balance/factory"
	walletDtos "finance-service/services/wallet/dto"
	"finance-service/services/wallet/transaction"
	transactionDtos "finance-service/services/wallet/transaction/dto"
	"finance-service/utils"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type DefaultWalletWriteService struct {
	BalanceHandlerFactory   *balanceHandlerFactory.BalanceHandlerFactory
	WalletRepository        *repositories.WalletRepository
	TransactionWriteService *transaction.TransactionWriteService
}

// TODO: Should we init all of these beans outside (for singleton) then pass it into the constructor here instead?
func NewWalletWriteServiceWithClient(
	balanceHandlerFactory *balanceHandlerFactory.BalanceHandlerFactory,
	walletRepository *repositories.WalletRepository,
	transactionWriteService *transaction.TransactionWriteService,
) *DefaultWalletWriteService {
	return &DefaultWalletWriteService{
		BalanceHandlerFactory:   balanceHandlerFactory,
		WalletRepository:        walletRepository,
		TransactionWriteService: transactionWriteService,
	}
}

func (this *DefaultWalletWriteService) UpdateBalance(tx *gorm.DB, updateRequest walletDtos.WalletUpdateRequest) (*walletDtos.WalletDto, error) {
	var result *walletDtos.WalletDto
	err := txManagement.WithTransaction(this.WalletRepository.DB, tx, func(localTx *gorm.DB) error {
		operationType, err := this.getOperationType(updateRequest.UpdateType)
		if err != nil {
			return err
		}

		walletId, err := this.getWalletId(updateRequest.ExternalWalletId)
		if err != nil {
			return err
		}

		if err := this.updateWalletBalance(localTx, walletId, updateRequest, operationType); err != nil {
			return err
		}

		transactionDto := transactionDtos.TransactionDto{
			WalletID: walletId,
			Amount:   updateRequest.Amount,
			Type:     updateRequest.UpdateType,
			Content:  updateRequest.Content,
		}

		if err := this.TransactionWriteService.CreateTransaction(localTx, transactionDto); err != nil {
			return err
		}

		result, err = this.getUpdatedWallet(localTx, walletId)
		return err
	})

	return result, err
}

func (this *DefaultWalletWriteService) WalletTransfer(walletTransferRequest walletDtos.WalletTransferRequest) ([]walletDtos.WalletDto, error) {
	amount, toExternalWalletId, fromExternalWalletId := walletTransferRequest.Amount, walletTransferRequest.ExternalToWalletId, walletTransferRequest.ExternalFromWalletId
	_, err := this.validateTransferAmount(amount)
	if err != nil {
		return nil, err
	}

	var updatedAsmWallet *walletDtos.WalletDto
	var updatedVndWallet *walletDtos.WalletDto
	// Begin new transaction with desired isolation level (REPEATABLE READ or SERIALIZABLE)
	err = txManagement.WithNewTransaction(this.WalletRepository.DB, sql.LevelRepeatableRead, func(tx *gorm.DB) error {
		_, err = this.validateWallets(tx, toExternalWalletId, fromExternalWalletId, amount)
		if err != nil {
			return err
		}

		// Update `from` wallet (debit)
		updatedVndWallet, err = this.UpdateBalance(tx, walletDtos.WalletUpdateRequest{
			ExternalWalletId: fromExternalWalletId,
			UpdateType:       operationTypes.VNDWalletDebit.String(),
			Amount:           amount,
			Content:          "Chuyen tien tu VND Wallet sang ASM Wallet",
		})
		if err != nil {
			return fmt.Errorf("failed to update from wallet: %v", err)
		}

		// Update `to` wallet (credit)
		updatedAsmWallet, err = this.UpdateBalance(tx, walletDtos.WalletUpdateRequest{
			ExternalWalletId: toExternalWalletId,
			UpdateType:       operationTypes.ASMWalletTopup.String(),
			Amount:           amount,
			Content:          "Nhan tien tu VND Wallet",
		})
		if err != nil {
			return fmt.Errorf("failed to update to wallet: %v", err)
		}

		utils.Logger().Println("Converted wallet balance", toExternalWalletId, amount)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Return both wallets
	return []walletDtos.WalletDto{*updatedVndWallet, *updatedAsmWallet}, nil
}

func (this *DefaultWalletWriteService) validateTransferAmount(amount float64) (*models.Wallet, error) {
	if amount < 2000000 {
		utils.Logger().Println("Transfer amount must be greater than 2.000.000")
		return nil, errors.New("transfer amount must be greater than 2.000.000")
	}
	return nil, nil
}

func (this *DefaultWalletWriteService) validateWallets(tx *gorm.DB, toExternalWalletId string, fromExternalWalletId string, amount float64) (*models.Wallet, error) {
	toWalletIdStr, _ := cryptography.Decrypt(toExternalWalletId)
	toWalletId, _ := strconv.ParseUint(toWalletIdStr, 10, 64)

	fromWalletIdStr, _ := cryptography.Decrypt(fromExternalWalletId)
	fromWalletId, _ := strconv.ParseUint(fromWalletIdStr, 10, 64)

	fromWallet, err := this.WalletRepository.GetByID(tx, uint(fromWalletId))
	if err != nil {
		utils.Logger().Println("Error getting wallet:", err)
		return nil, errors.New("wallet not found")
	}

	_, err = this.WalletRepository.GetByID(tx, uint(toWalletId))
	if err != nil {
		utils.Logger().Println("Error getting wallet:", err)
		return nil, errors.New("wallet not found")
	}

	if fromWallet.Balance < amount {
		utils.Logger().Println("Insufficient wallet balance")
		return nil, errors.New("insufficient wallet balance")
	}
	return nil, nil
}

func (this *DefaultWalletWriteService) getOperationType(topupType string) (operationTypes.BalanceOperation, error) {
	operationType, err := operationTypes.ParseTopupType(topupType)
	if err != nil {
		utils.Logger().Println("Error parsing topup type:", err)
		return operationType, errors.New("invalid topup type")
	}
	return operationType, nil
}

func (this *DefaultWalletWriteService) getWalletId(encryptedWalletId string) (uint, error) {
	walletIdStr, err := cryptography.Decrypt(encryptedWalletId)
	if err != nil {
		return 0, fmt.Errorf("error decrypting wallet ID: %v", err)
	}

	walletIdLong, err := strconv.ParseUint(walletIdStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting wallet ID to uint: %v", err)
	}

	return uint(walletIdLong), nil
}

func (this *DefaultWalletWriteService) updateWalletBalance(tx *gorm.DB, walletId uint, updateRequest walletDtos.WalletUpdateRequest, operationType operationTypes.BalanceOperation) error {
	balanceHandler := this.BalanceHandlerFactory.GetHandler(operationType)
	if balanceHandler == nil {
		utils.Logger().Println("Operation type not found:", operationType)
		return errors.New("operation type not found")
	}

	err := balanceHandler.UpdateBalance(tx, balanceDtos.UpdateBalanceInput{
		WalletId:         walletId,
		WalletType:       updateRequest.WalletType,
		Amount:           updateRequest.Amount,
		BalanceOperation: operationType,
	})
	if err != nil {
		return err
	}

	return nil
}

func (this *DefaultWalletWriteService) getUpdatedWallet(tx *gorm.DB, walletId uint) (*walletDtos.WalletDto, error) {
	var wallet models.Wallet
	if err := tx.First(&wallet, walletId).Error; err != nil {
		return nil, err
	}

	utils.Logger().Println("Wallet balance updated", walletId)
	return &walletDtos.WalletDto{
		Balance:    wallet.Balance,
		UserId:     wallet.UserId,
		WalletType: wallet.WalletType,
	}, nil
}

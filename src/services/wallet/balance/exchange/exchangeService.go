package services

import (
	"context"
	"errors"
	"finance-service/models"
	"finance-service/repositories"
	userrpcclient "finance-service/rpc/client"
	walletservices "finance-service/services/wallet"
	"finance-service/utils"
	"time"
)

// ExchangeService provides operations for converting balances.
type ExchangeService struct {
	WalletService         *walletservices.DefaultWalletWriteService
	UserServiceClient     *userrpcclient.UserServiceClient
	TransactionRepository *repositories.TransactionRepository
}

// NewExchangeService creates a new ExchangeService.
func NewExchangeService(walletService *walletservices.DefaultWalletWriteService, userServiceClient *userrpcclient.UserServiceClient, transactionRepository *repositories.TransactionRepository) *ExchangeService {
	return &ExchangeService{
		WalletService:         walletService,
		UserServiceClient:     userServiceClient,
		TransactionRepository: transactionRepository,
	}
}

// ConvertBalance converts the balance from one type to another.
func (service *ExchangeService) ConvertBalance(ctx context.Context, uuid string, amount float64) (*models.Wallet, error) {
	if amount < 2000000 {
		utils.Logger().Println("Transfer amount must be greater than 2.000.000")
		return nil, errors.New("transfer amount must be greater than 2.000.000")
	}

	// Fetch user data from user service
	user, err := service.UserServiceClient.GetUser(uuid)
	if err != nil {
		utils.Logger().Println("Error finding user:", err)
		return nil, errors.New("user is not registered")
	}

	// Fetch wallet data
	wallet, err := service.WalletService.WalletRepository.GetByUserID(user.Id)
	if err != nil {
		utils.Logger().Println("Error getting wallet:", err)
		return nil, errors.New("wallet not found")
	}

	if wallet.Balance < amount {
		utils.Logger().Println("Insufficient wallet balance")
		return nil, errors.New("insufficient wallet balance")
	}

	newBalance := wallet.Balance - amount

	if err := service.WalletService.WalletRepository.UpdateBalance(wallet.ID, newBalance); err != nil {
		utils.Logger().Println("Error updating wallet balance:", err)
		return nil, err
	}

	transaction1 := models.Transaction{
		WalletID:  wallet.ID,
		Amount:    amount,
		Type:      "debit",
		Content:   "Convert balance debit",
		CreatedBy: "system",
		CreatedOn: time.Now(),
	}
	if err := service.TransactionRepository.Create(transaction1); err != nil {
		utils.Logger().Println("Error creating transaction:", err)
		return nil, err
	}

	transaction2 := models.Transaction{
		WalletID:  wallet.ID,
		Amount:    amount,
		Type:      "credit",
		Content:   "Convert balance credit",
		CreatedBy: "system",
		CreatedOn: time.Now(),
	}
	if err := service.TransactionRepository.Create(transaction2); err != nil {
		utils.Logger().Println("Error creating transaction:", err)
		return nil, err
	}

	utils.Logger().Println("Converted wallet balance", wallet.ID, amount)
	return wallet, nil
}

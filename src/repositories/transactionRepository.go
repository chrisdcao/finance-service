package repositories

import (
	"finance-service/config"
	"finance-service/models"
)

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) Create(transaction models.Transaction) error {
	return config.DB.Create(&transaction).Error
}

func (r *TransactionRepository) GetAllByWalletID(walletID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := config.DB.Where("wallet_id = ?", walletID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetAllByUuid(uuid string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := config.DB.Where("uuid = ?", uuid).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

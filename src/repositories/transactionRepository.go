package repositories

import (
	"finance-service/models"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) Create(tx *gorm.DB, transaction *models.Transaction) error {
	return tx.Create(transaction).Error
}

func (r *TransactionRepository) GetAllByWalletID(walletID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.DB.Where("wallet_id = ?", walletID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) GetAllByUuid(uuid string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.DB.Where("uuid = ?", uuid).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

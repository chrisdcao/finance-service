package repositories

import (
	"finance-service/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type TransactionRepository struct {
	DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) Create(tx *gorm.DB, transaction models.Transaction) (*models.Transaction, error) {
	err := tx.Create(&transaction).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to create transaction")
	}
	return &transaction, nil
}

func (r *TransactionRepository) FindTransactions(walletType, actionType string, amount float64, fromTime, toTime time.Time, uuid string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.DB.Model(&models.Transaction{})

	if walletType != "" {
		query = query.Where("wallet_type = ?", walletType)
	}
	if actionType != "" {
		query = query.Where("type = ?", actionType)
	}
	if amount != 0 {
		query = query.Where("amount = ?", amount)
	}
	if !fromTime.IsZero() {
		query = query.Where("created_on >= ?", fromTime)
	}
	if !toTime.IsZero() {
		query = query.Where("created_on <= ?", toTime)
	}
	if uuid != "" {
		query = query.Where("uuid = ?", uuid)
	}

	err := query.Find(&transactions).Error
	return transactions, errors.Wrap(err, "failed to find transactions")
}

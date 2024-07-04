package repositories

import (
	"finance-service/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type WalletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{DB: db}
}

func (r *WalletRepository) UpdateBalance(tx *gorm.DB, walletID uint, newBalance float64) error {
	err := tx.Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", newBalance).Error
	if err != nil {
		return errors.Wrap(err, "failed to update wallet balance")
	}
	return nil
}

func (r *WalletRepository) GetByID(tx *gorm.DB, walletID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := tx.First(&wallet, walletID).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get wallet")
	}
	return &wallet, nil
}

func (r *WalletRepository) GetByUserID(tx *gorm.DB, userID uint) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get wallet")
	}
	return &wallet, nil
}

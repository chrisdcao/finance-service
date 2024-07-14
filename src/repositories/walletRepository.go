package repositories

import (
	"finance-service/models"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type WalletRepository struct {
	DB *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{DB: db}
}

func (r *WalletRepository) UpdateBalance(tx *gorm.DB, wallet models.Wallet) (*models.Wallet, error) {
	err := tx.Model(&wallet).Where("id = ?", wallet.Id).Update("balance", wallet.Balance).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to update wallet balance")
	}
	return &wallet, nil
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

func (this *WalletRepository) GetByUserIDAndWalletType(tx *gorm.DB, userID string, walletType string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := tx.Where("user_id = ? AND wallet_type = ?", userID, walletType).First(&wallet).Error; err != nil {
		errMsg := fmt.Sprintf("failed to get wallet with [user_id: %s] and [wallet_type: %s]", userID, walletType)
		return nil, errors.Wrap(err, errMsg)
	}
	return &wallet, nil
}

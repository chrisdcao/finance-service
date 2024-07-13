package services

import (
	"context"
	"finance-service/repositories"
	walletDtos "finance-service/services/wallet/dto"
	walletservices "finance-service/services/wallet/mapper"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DefaultWalletReadService struct {
	WalletRepository *repositories.WalletRepository
	WalletMapper     *walletservices.WalletMapper
}

func NewWalletReadService(repository *repositories.WalletRepository, walletMapper *walletservices.WalletMapper) *DefaultWalletReadService {
	return &DefaultWalletReadService{
		WalletRepository: repository,
		WalletMapper:     walletMapper,
	}
}

func (this *DefaultWalletReadService) GetWalletByUserId(ctx context.Context, tx *gorm.DB, userId uint) (*walletDtos.WalletDto, error) {
	wallet, err := this.WalletRepository.GetByUserID(tx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wallet")
	}

	return this.WalletMapper.FromModelToDto(*wallet), nil
}

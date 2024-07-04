package services

import (
	"context"
	"finance-service/models"
	"finance-service/repositories"
	walletDtos "finance-service/services/wallet/dto"
	"finance-service/services/wallet/parser"
	"gorm.io/gorm"
)

type DefaultWalletReadService struct {
	WalletRepository *repositories.WalletRepository
	WalletIdParser   *parser.WalletIdParser
}

func NewWalletReadService(repository *repositories.WalletRepository, idParser *parser.WalletIdParser) *DefaultWalletReadService {
	return &DefaultWalletReadService{
		WalletRepository: repository,
		WalletIdParser:   idParser,
	}
}

func (this *DefaultWalletReadService) GetFromExternalId(ctx context.Context, externalWalletId string, tx *gorm.DB) (*models.Wallet, error) {
	walletId, err := this.WalletIdParser.ParseFromEncryption(ctx, externalWalletId)
	if err != nil {
		return nil, err
	}

	wallet, err := this.WalletRepository.GetByID(tx, walletId)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (this *DefaultWalletReadService) GetWallet(ctx context.Context, tx *gorm.DB, walletId uint) (*walletDtos.WalletDto, error) {
	wallet, err := this.WalletRepository.GetByID(tx, walletId)
	if err != nil {
		return nil, err
	}

	return &walletDtos.WalletDto{
		Balance:    wallet.Balance,
		UserId:     wallet.UserId,
		WalletType: wallet.WalletType,
	}, nil
}

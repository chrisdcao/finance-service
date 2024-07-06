package services

import (
	"context"
	"finance-service/models"
	"finance-service/repositories"
	walletDtos "finance-service/services/wallet/dto"
	"finance-service/services/wallet/parser"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "failed to parse wallet id")
	}

	wallet, err := this.WalletRepository.GetByID(tx, walletId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wallet")
	}

	return wallet, nil
}

func (this *DefaultWalletReadService) GetWallet(ctx context.Context, tx *gorm.DB, walletId uint) (*walletDtos.WalletDto, error) {
	wallet, err := this.WalletRepository.GetByID(tx, walletId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wallet")
	}

	walletDto := &walletDtos.WalletDto{
		Balance:    wallet.Balance,
		UserId:     wallet.UserId,
		WalletType: wallet.WalletType,
	}

	return walletDto, nil
}

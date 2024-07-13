package read

import (
	"context"
	"finance-service/repositories"
	walletDtos "finance-service/services/wallet/dto"
	walletservices "finance-service/services/wallet/mapper"
	"finance-service/services/wallet/parser"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DefaultWalletReadService struct {
	WalletRepository *repositories.WalletRepository
	WalletIdParser   *parser.WalletIdParser
	WalletMapper     *walletservices.WalletMapper
}

func NewWalletReadService(repository *repositories.WalletRepository, idParser *parser.WalletIdParser, walletMapper *walletservices.WalletMapper) *DefaultWalletReadService {
	return &DefaultWalletReadService{
		WalletRepository: repository,
		WalletIdParser:   idParser,
		WalletMapper:     walletMapper,
	}
}

func (this *DefaultWalletReadService) GetFromExternalId(ctx context.Context, externalWalletId string, tx *gorm.DB) (*walletDtos.WalletDto, error) {
	walletId, err := this.WalletIdParser.ParseFromEncryption(ctx, externalWalletId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse wallet id")
	}

	wallet, err := this.WalletRepository.GetByID(tx, walletId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wallet")
	}

	return this.WalletMapper.FromModelToDto(*wallet), nil
}

func (this *DefaultWalletReadService) GetWalletByUserId(ctx context.Context, tx *gorm.DB, userId uint) (*walletDtos.WalletDto, error) {
	wallet, err := this.WalletRepository.GetByUserID(tx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wallet")
	}

	return this.WalletMapper.FromModelToDto(*wallet), nil
}

func (this *DefaultWalletReadService) GetWalletById(ctx context.Context, tx *gorm.DB, walletId uint) (*walletDtos.WalletDto, error) {
	wallet, err := this.WalletRepository.GetByID(tx, walletId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get wallet")
	}

	return this.WalletMapper.FromModelToDto(*wallet), nil
}

package mapper

import (
	"finance-service/models"
	walletDtos "finance-service/services/wallet/dto"
)

type WalletMapper struct{}

func NewWalletMapper() *WalletMapper {
	return &WalletMapper{}
}

func (this *WalletMapper) FromDtoToModel(dto walletDtos.WalletDto) *models.Wallet {
	return &models.Wallet{
		Balance:    dto.Balance,
		UserId:     dto.UserId,
		WalletType: dto.WalletType,
	}
}

func (this *WalletMapper) FromModelToDto(wallet models.Wallet) *walletDtos.WalletDto {
	return &walletDtos.WalletDto{
		Balance:    wallet.Balance,
		UserId:     wallet.UserId,
		WalletType: wallet.WalletType,
	}
}

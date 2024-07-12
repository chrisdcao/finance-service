package services

import (
	"finance-service/models"
	transactionDtos "finance-service/services/balance/dto"
	walletDtos "finance-service/services/wallet/dto"
)

type WalletMapper struct{}

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

func (this *WalletMapper) FromUpdateBalanceInputToWalletModel(input transactionDtos.UpdateBalanceInput) (*models.Wallet, error) {
	return nil, nil
}

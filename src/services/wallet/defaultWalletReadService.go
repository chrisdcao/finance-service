package services

import (
	"finance-service/repositories"
	walletservices "finance-service/services/wallet/mapper"
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

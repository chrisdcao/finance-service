package services

import (
	"finance-service/controllers/dto/request"
	"finance-service/models"
	walletDtos "finance-service/services/wallet/dto"
)

type WalletWriteService interface {
	UpdateBalance(topupRequest request.WalletUpdateRequest) (*walletDtos.WalletDto, error)
	ConvertBalance(toExternalWalletId string, fromExternalWalletId string, amount float64) (*models.Wallet, error)
}

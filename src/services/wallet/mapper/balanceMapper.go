package mapper

import (
	walletDtos "finance-service/controllers/wallet/dto/request"
	"finance-service/services/wallet/dto"
	operationTypes "finance-service/services/wallet/enums"
	walletTypes "finance-service/services/wallet/enums"
	"github.com/pkg/errors"
)

type BalanceMapper struct{}

func NewBalanceMapper() *BalanceMapper {
	return &BalanceMapper{}
}

func (this *BalanceMapper) FromUpdateRequestToInput(request walletDtos.WalletUpdateRequest) (*dto.UpdateBalanceInput, error) {
	operationType, err := operationTypes.ParseStringToBalanceOperation(request.UpdateType)
	if err != nil {
		return nil, errors.Wrap(err, "Unrecognized operation type: "+request.UpdateType)
	}

	walletType, err := walletTypes.ParseStringToWalletType(request.WalletType)
	if err != nil {
		return nil, errors.Wrap(err, "Unrecognized wallet type: "+request.WalletType)
	}

	return &dto.UpdateBalanceInput{
		UserId:           request.UserId,
		WalletType:       walletType,
		DiffAmount:       request.Amount,
		BalanceOperation: operationType,
		Content:          request.Content,
	}, nil
}

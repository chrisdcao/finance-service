package mapper

import (
	walletDtos "finance-service/controllers/dto/request"
	"finance-service/services/balance/dto"
	operationTypes "finance-service/services/balance/enums"
	"github.com/pkg/errors"
)

type BalanceMapper struct {
}

func NewBalanceMapper() *BalanceMapper {
	return &BalanceMapper{}
}

func (this *BalanceMapper) FromUpdateRequestToInput(request walletDtos.WalletUpdateRequest) (*dto.UpdateBalanceInput, error) {
	operationType, err := operationTypes.Parse(request.UpdateType)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert to operationType: "+request.UpdateType)
	}

	return &dto.UpdateBalanceInput{
		UserId:           request.UserId,
		WalletType:       request.WalletType,
		DiffAmount:       request.Amount,
		BalanceOperation: operationType,
		Content:          request.Content,
	}, nil
}

package factory

import (
	"finance-service/repositories"
	transactionDtos "finance-service/services/wallet/dto"
	"finance-service/services/wallet/enums"
	balanceTypes "finance-service/services/wallet/enums"
	"finance-service/services/wallet/handler"
	"finance-service/services/wallet/handler/credit"
	"finance-service/services/wallet/handler/debit"
	walletservices "finance-service/services/wallet/mapper"
	"finance-service/services/wallet/validator"
	"github.com/pkg/errors"
)

type BalanceHandlerFactory struct {
	handlers map[enums.BalanceOperation]handler.BalanceHandler
}

// NewBalanceHandlerFactory returns a new instance of BalanceHandlerFactory with ds: <balanceType, handler>
func NewBalanceHandlerFactory(walletRepo *repositories.WalletRepository, walletValidator *validator.DefaultWalletValidator, walletMapper *walletservices.WalletMapper) *BalanceHandlerFactory {
	factory := &BalanceHandlerFactory{
		handlers: make(map[enums.BalanceOperation]handler.BalanceHandler),
	}

	// Initialize handlers
	creditHandler := credit.NewCreditBalanceHandler(walletRepo, walletMapper)
	debitHandler := debit.NewDebitBalanceHandler(walletRepo, walletMapper, walletValidator)

	// Register handlers
	factory.RegisterHandler(balanceTypes.Credit, creditHandler)
	factory.RegisterHandler(balanceTypes.Debit, debitHandler)

	return factory
}

func (this *BalanceHandlerFactory) RegisterHandler(transactionType enums.BalanceOperation, handler handler.BalanceHandler) {
	this.handlers[transactionType] = handler
}

func (this *BalanceHandlerFactory) GetHandler(updateInput transactionDtos.UpdateBalanceInput) (handler.BalanceHandler, error) {
	hand, exists := this.handlers[updateInput.BalanceOperation]
	if !exists {
		return nil, errors.New("handler for operation type" + updateInput.BalanceOperation.String() + " not found")
	}
	return hand, nil
}

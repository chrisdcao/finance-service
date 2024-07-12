package factory

import (
	"finance-service/repositories"
	transactionDtos "finance-service/services/balance/dto"
	balanceTypes "finance-service/services/balance/enums"
	"finance-service/services/balance/handler"
	"finance-service/services/balance/handler/credit"
	"finance-service/services/balance/handler/debit"
	"github.com/pkg/errors"
)

type BalanceHandlerFactory struct {
	handlers map[balanceTypes.BalanceOperation]handler.BalanceHandler
}

// NewBalanceHandlerFactory returns a new instance of BalanceHandlerFactory with ds: <balanceType, handler>
func NewBalanceHandlerFactory(walletRepo *repositories.WalletRepository) *BalanceHandlerFactory {
	factory := &BalanceHandlerFactory{
		handlers: make(map[balanceTypes.BalanceOperation]handler.BalanceHandler),
	}

	// Initialize handlers
	// TODO: Ask Huys and add concrete impl of missing handlers
	asmWalletDebitHandler := debit.NewDebitTransaction(walletRepo)
	asmWalletTopupHandler := credit.NewCreditTransaction(walletRepo)
	vndWalletDebitHandler := debit.NewDebitTransaction(walletRepo)
	vndWalletTopupHandler := credit.NewCreditTransaction(walletRepo)

	// Register handlers
	factory.RegisterHandler(balanceTypes.ASMWalletDebit, asmWalletDebitHandler)
	factory.RegisterHandler(balanceTypes.ASMWalletTopup, asmWalletTopupHandler)
	factory.RegisterHandler(balanceTypes.VNDWalletDebit, vndWalletDebitHandler)
	factory.RegisterHandler(balanceTypes.VNDWalletTopup, vndWalletTopupHandler)

	return factory
}

func (this *BalanceHandlerFactory) RegisterHandler(transactionType balanceTypes.BalanceOperation, handler handler.BalanceHandler) {
	this.handlers[transactionType] = handler
}

func (this *BalanceHandlerFactory) GetHandler(updateInput transactionDtos.UpdateBalanceInput) (handler.BalanceHandler, error) {
	hand, exists := this.handlers[updateInput.BalanceOperation]
	if !exists {
		return nil, errors.New("handler for operation type" + updateInput.BalanceOperation.String() + " not found")
	}
	return hand, nil
}

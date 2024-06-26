package factory

import (
	"finance-service/repositories"
	balanceTypes "finance-service/services/wallet/balance/enums"
	"finance-service/services/wallet/balance/handler"
	"finance-service/services/wallet/balance/handler/credit"
	"finance-service/services/wallet/balance/handler/debit"
)

type BalanceHandlerFactory struct {
	handlers map[balanceTypes.BalanceOperation]handler.BalanceHandler
}

// Setting up the handlers and the map
func InitializeTransactionHandlers(repo *repositories.WalletRepository) *BalanceHandlerFactory {
	factory := NewBalanceHandlerFactory()

	// Initialize handlers
	asmWalletDebitHandler := debit.NewDebitTransaction(repo)
	asmWalletTopupHandler := credit.NewCreditTransaction(repo)
	vndWalletDebitHandler := debit.NewDebitTransaction(repo)
	vndWalletTopupHandler := credit.NewCreditTransaction(repo)

	// Register handlers
	factory.RegisterHandler(balanceTypes.ASMWalletDebit, asmWalletDebitHandler)
	factory.RegisterHandler(balanceTypes.ASMWalletTopup, asmWalletTopupHandler)
	factory.RegisterHandler(balanceTypes.VNDWalletDebit, vndWalletDebitHandler)
	factory.RegisterHandler(balanceTypes.VNDWalletTopup, vndWalletTopupHandler)

	return factory
}

func NewBalanceHandlerFactory() *BalanceHandlerFactory {
	return &BalanceHandlerFactory{
		handlers: make(map[balanceTypes.BalanceOperation]handler.BalanceHandler),
	}
}

func (f *BalanceHandlerFactory) RegisterHandler(transactionType balanceTypes.BalanceOperation, handler handler.BalanceHandler) {
	f.handlers[transactionType] = handler
}

func (f *BalanceHandlerFactory) GetHandler(txType balanceTypes.BalanceOperation) handler.BalanceHandler {
	handler, exists := f.handlers[txType]
	if !exists {
		return nil
	}
	return handler
}

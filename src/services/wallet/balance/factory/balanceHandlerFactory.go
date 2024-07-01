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

func (this *BalanceHandlerFactory) GetHandler(txType balanceTypes.BalanceOperation) handler.BalanceHandler {
	hand, exists := this.handlers[txType]
	if !exists {
		return nil
	}
	return hand
}

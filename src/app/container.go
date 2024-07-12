package app

import (
	"finance-service/config"
	"finance-service/controllers"
	"finance-service/repositories"
	"finance-service/services/balance/factory"
	transaction2 "finance-service/services/transaction"
	"finance-service/services/transaction/mapper"
	walletservices "finance-service/services/wallet"
	"finance-service/services/wallet/parser"
	"finance-service/services/wallet/read"
	"finance-service/services/wallet/validator"
	"finance-service/services/wallet/write"
)

type Container struct {
	EndUserController *controllers.EndUserController
	AdminController   *controllers.AdminController
}

func NewContainer() *Container {
	// Initialize services
	transactionRepository := repositories.NewTransactionRepository(config.DB)
	transactionReadService := transaction2.NewTransactionReadService(transactionRepository)
	transactionWriteService := transaction2.NewTransactionWriteService(transactionRepository)
	transactionMapper := mapper.NewTransactionMapper()

	// Initialize Repos
	walletRepository := repositories.NewWalletRepository(config.DB)
	balanceHandlerFactory := factory.NewBalanceHandlerFactory(walletRepository)

	walletIdParser := parser.NewWalletIdParser()
	walletReadService := read.NewWalletReadService(walletRepository, walletIdParser)
	walletValidator := validator.NewWalletValidator(walletReadService)

	walletWriteService := write.NewWalletWriteService(
		balanceHandlerFactory,
		walletRepository,
		transactionWriteService,
		transactionMapper,
		walletValidator,
		walletIdParser,
	)

	walletService := walletservices.NewWalletService(
		balanceHandlerFactory,
		transactionWriteService,
		transactionMapper,
		walletValidator,
		walletReadService,
		walletWriteService,
	)

	// Initialize Controllers
	endUserController := controllers.NewEndUserController(transactionReadService)
	adminController := controllers.NewAdminController(walletService)

	return &Container{
		EndUserController: endUserController,
		AdminController:   adminController,
	}
}

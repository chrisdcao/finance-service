package di

import (
	"finance-service/configs/db/connection"
	"finance-service/controllers"
	"finance-service/repositories"
	transaction2 "finance-service/services/transaction"
	"finance-service/services/transaction/mapper"
	walletservices "finance-service/services/wallet"
	"finance-service/services/wallet/factory"
	mapper2 "finance-service/services/wallet/mapper"
	"finance-service/services/wallet/validator"
)

type Container struct {
	EndUserController *controllers.EndUserController
	AdminController   *controllers.AdminController
}

func NewContainer() *Container {
	// Initialize services
	transactionRepository := repositories.NewTransactionRepository(connection.DB)
	transactionReadService := transaction2.NewTransactionReadService(transactionRepository)
	transactionWriteService := transaction2.NewTransactionWriteService(transactionRepository)
	transactionMapper := mapper.NewTransactionMapper()

	// Initialize Repos
	walletRepository := repositories.NewWalletRepository(connection.DB)

	// init mapper
	balanceMapper := mapper2.NewBalanceMapper()
	walletMapper := mapper2.NewWalletMapper()

	// init services
	walletReadService := walletservices.NewWalletReadService(walletRepository, walletMapper)
	walletValidator := validator.NewWalletValidator()
	balanceHandlerFactory := factory.NewBalanceHandlerFactory(walletRepository, walletValidator, walletMapper)

	walletWriteService := walletservices.NewWalletWriteService(balanceHandlerFactory, walletRepository, transactionWriteService, transactionMapper, balanceMapper, walletValidator)

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

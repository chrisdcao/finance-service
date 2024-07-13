package configs

import (
	"finance-service/controllers"
	"finance-service/repositories"
	transaction2 "finance-service/services/transaction"
	"finance-service/services/transaction/mapper"
	walletservices "finance-service/services/wallet"
	"finance-service/services/wallet/factory"
	"finance-service/services/wallet/parser"
	"finance-service/services/wallet/validator"
)

type Container struct {
	EndUserController *controllers.EndUserController
	AdminController   *controllers.AdminController
}

func NewContainer() *Container {
	// Initialize services
	transactionRepository := repositories.NewTransactionRepository(DB)
	transactionReadService := transaction2.NewTransactionReadService(transactionRepository)
	transactionWriteService := transaction2.NewTransactionWriteService(transactionRepository)
	transactionMapper := mapper.NewTransactionMapper()

	// Initialize Repos
	walletRepository := repositories.NewWalletRepository(DB)
	balanceHandlerFactory := factory.NewBalanceHandlerFactory(walletRepository)

	walletIdParser := parser.NewWalletIdParser()
	walletReadService := walletservices.NewWalletReadService(walletRepository, walletIdParser)
	walletValidator := validator.NewWalletValidator(walletReadService)

	walletWriteService := walletservices.NewWalletWriteService(
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

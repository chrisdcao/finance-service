package app

import (
	"finance-service/config"
	"finance-service/controllers"
	"finance-service/repositories"
	walletservices "finance-service/services/wallet"
	"finance-service/services/wallet/balance/factory"
	"finance-service/services/wallet/parser"
	"finance-service/services/wallet/transaction"
	"finance-service/services/wallet/transaction/mapper"
	"finance-service/services/wallet/validator"
)

//const userServiceAddress = "localhost:50051" // Address of the user service

type Container struct {
	EndUserController *controllers.EndUserController
	AdminController   *controllers.AdminController
}

func NewContainer() *Container {
	// Initialize gRPC client for user service (NOT USED FOR NOW)
	//_ = setupUserServiceClient()

	// Initialize services
	transactionRepository := repositories.NewTransactionRepository(config.DB)
	transactionReadService := transaction.NewTransactionReadService(transactionRepository)
	transactionWriteService := transaction.NewTransactionWriteService(transactionRepository)
	transactionMapper := mapper.NewTransactionMapper()

	// Initialize Repos
	walletRepository := repositories.NewWalletRepository(config.DB)
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

//func setupUserServiceClient() *userrpcclient.UserServiceClient {
//	conn, err := grpc.Dial(userServiceAddress, grpc.WithInsecure(), grpc.WithBlock())
//	if err != nil {
//		log.Fatalf("Failed to connect to user service: %v", err)
//	}
//	userServiceClient := generated_files.NewUserServiceClient(conn)
//	return userrpcclient.NewUserServiceClient(userServiceClient)
//}

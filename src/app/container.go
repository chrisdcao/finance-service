package app

import (
	"finance-service/config"
	"finance-service/controllers"
	"finance-service/repositories"
	userrpcclient "finance-service/rpc/client"
	"finance-service/rpc/protos/user/generated_files"
	walletservices "finance-service/services/wallet"
	"finance-service/services/wallet/balance/factory"
	"finance-service/services/wallet/transaction"
	"google.golang.org/grpc"
	"log"
)

// TODO: Refactor this to .env file
const userServiceAddress = "localhost:50051" // Address of the user service

type Container struct {
	EndUserController *controllers.EndUserController
	AdminController   *controllers.AdminController
}

// TODO: Clean up `main` once this `Container` already has all the Beans initialization
func NewContainer() *Container {
	// Initialize gRPC client for user service (NOT USED FOR NOW)
	_ = setupUserServiceClient()

	// Initialize Repos
	walletRepository := repositories.NewWalletRepository(config.DB)

	// Initialize services
	transactionRepository := repositories.NewTransactionRepository(config.DB)
	transactionReadService := transaction.NewTransactionReadService(transactionRepository)
	transactionWriteService := transaction.NewTransactionWriteService(transactionRepository)

	balanceHandlerFactory := factory.NewBalanceHandlerFactory(walletRepository)

	walletWriteService := walletservices.NewWalletService(
		balanceHandlerFactory,
		walletRepository,
		transactionWriteService,
	)

	// Initialize Controllers
	endUserController := controllers.NewEndUserController(transactionReadService)
	adminController := controllers.NewAdminController(walletWriteService)

	return &Container{
		EndUserController: endUserController,
		AdminController:   adminController,
	}
}

func setupUserServiceClient() *userrpcclient.UserServiceClient {
	conn, err := grpc.Dial(userServiceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	userServiceClient := generated_files.NewUserServiceClient(conn)
	return userrpcclient.NewUserServiceClient(userServiceClient)
}

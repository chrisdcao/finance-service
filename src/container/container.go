package container

import (
	"finance-service/controllers"
	"finance-service/repositories"
	userrpcclient "finance-service/rpc/client"
	"finance-service/rpc/protos/user/generated_files"
	walletservices "finance-service/services/wallet"
	exchangeservices "finance-service/services/wallet/balance/exchange"
	"google.golang.org/grpc"
	"log"
)

const userServiceAddress = "localhost:50051" // Address of the user service

type Container struct {
	TransactionController *controllers.TransactionController
	WalletController      *controllers.WalletController
}

func NewContainer() *Container {
	// Initialize gRPC client for user service
	userServiceClientWrapper := setupUserServiceClient()

	// Initialize repositories
	transactionRepository := repositories.NewTransactionRepository()

	// Initialize services
	walletService := walletservices.NewWalletServiceWithClient(userServiceClientWrapper)
	exchangeService := exchangeservices.NewExchangeService(walletService, userServiceClientWrapper, transactionRepository)

	// Initialize controllers
	transactionController := controllers.NewTransactionController()
	walletController := controllers.NewWalletController(walletService, exchangeService)

	return &Container{
		TransactionController: transactionController,
		WalletController:      walletController,
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

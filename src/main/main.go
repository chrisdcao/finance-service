package main

import (
	"context"
	"finance-service/config"
	"finance-service/models"
	"finance-service/repositories"
	userrpcclient "finance-service/rpc/client"
	"finance-service/rpc/protos/user/generated_files"
	walletservices "finance-service/services/wallet"
	exchangeservices "finance-service/services/wallet/balance/exchange"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	userServiceAddress = "localhost:50051" // Address of the user service
	grpcPort           = ":50052"          // Port for the finance gRPC service
	httpPort           = ":3000"           // Port for the HTTP server
)

func main() {
	// Load environment variables
	loadEnv()

	// Initialize the database
	initDatabase()

	// Auto migrate the models
	autoMigrateModels()

	// Set up routes
	router := setupRouter()

	// Set up gRPC client for user service
	userServiceClientWrapper := setupUserServiceClient()

	// Initialize services and repositories
	walletService, transactionRepository, exchangeService := initializeServices(userServiceClientWrapper)

	// Start HTTP server
	httpServer := startHTTPServer(router)

	// Graceful shutdown
	gracefulShutdown(httpServer)
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func initDatabase() {
	config.InitDB()
}

func autoMigrateModels() {
	config.DB.AutoMigrate(&models.Wallet{}, &models.Transaction{})
}

func setupUserServiceClient() *userrpcclient.UserServiceClient {
	conn, err := grpc.Dial(userServiceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	userServiceClient := generated_files.NewUserServiceClient(conn)
	return userrpcclient.NewUserServiceClient(userServiceClient)
}

func initializeServices(userServiceClientWrapper *userrpcclient.UserServiceClient) (*walletservices.DefaultWalletWriteService, *repositories.TransactionRepository, *exchangeservices.ExchangeService) {
	walletService := walletservices.NewWalletServiceWithClient(userServiceClientWrapper)
	transactionRepository := repositories.NewTransactionRepository()
	exchangeService := exchangeservices.NewExchangeService(walletService, userServiceClientWrapper, transactionRepository)
	return walletService, transactionRepository, exchangeService
}

func startHTTPServer(router http.Handler) *http.Server {
	server := &http.Server{
		Addr:    httpPort,
		Handler: router,
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Printf("Starting HTTP server on %s", httpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP server: %v", err)
		}
	}()

	return server
}

func gracefulShutdown(httpServer *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}
	log.Println("HTTP server stopped")
}

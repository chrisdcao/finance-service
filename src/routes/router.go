package routes

import (
	log2 "finance-service/common/log"
	"finance-service/configs/di"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	logger := log2.NewCommonLogger()
	router.Use(log2.TraceIDMiddleware(logger))

	// Setup DI container
	container := di.NewContainer()

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define transaction routes
	// TODO: based on the standard of REST APIs (using same endpoint, only with different methods for different operations)
	router.GET("/transactions/get", container.EndUserController.GetTransactions)

	// Define wallet routes
	router.POST("/wallets/updateBalance", container.AdminController.Topup)
	router.POST("/wallets/convertBalance", container.AdminController.WalletTransfer)

	return router
}

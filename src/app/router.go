package app

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// TODO: Move the constructed service to this place and have router using them?
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Setup DI container
	container := NewContainer()

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define transaction routes
	// TODO: based on the standard of REST APIs (using same endpoint with diff methods for diff interactions)
	router.GET("/transactions", container.EndUserController.GetTransactions)

	// Define wallet routes
	router.POST("/wallets/update_balance", container.AdminController.Topup)
	router.POST("/wallets/convert_balance", container.AdminController.WalletTransfer)

	return router
}

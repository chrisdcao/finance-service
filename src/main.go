package main

import (
	"finance-service/configs"
	_ "finance-service/controllers/docs" // Import for documents to be made
	"finance-service/routes"
)

// @title Finance Service API
// @version 1.0
// @description This is a sample finance service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /

func main() {
	// Load environment variables
	configs.LoadEnv()

	// Initialize the database
	configs.InitDatabase()

	// Auto migrate the models
	configs.AutoMigrateModels()

	// Set up routes
	router := routes.SetupRouter()

	// Start HTTP server
	httpServer := configs.StartHTTPServer(router)

	// Graceful shutdown
	configs.GracefulShutdown(httpServer)
}

package src

import (
	"finance-service/app"
	_ "finance-service/docs" // Import for documents to be made
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
	app.LoadEnv()

	// Initialize the database
	app.InitDatabase()

	// Auto migrate the models
	app.AutoMigrateModels()

	// Set up routes
	router := app.SetupRouter()

	// Start HTTP server
	httpServer := app.StartHTTPServer(router)

	// Graceful shutdown
	app.GracefulShutdown(httpServer)
}

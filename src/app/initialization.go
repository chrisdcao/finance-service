package app

import (
	"context"
	"finance-service/config"
	"finance-service/models"
	log2 "finance-service/utils/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	// TODO: Refactor this to .env file
	httpPort = ":3000" // Port for the HTTP server
)

func LoadEnv() log2.FunctionWithNoArgs {
	return log2.LogWraps(LoadEnvWithoutLogs)
}

func InitDatabase() log2.FunctionWithNoArgs {
	return log2.LogWraps(InitDatabaseWithoutLogs)
}

func AutoMigrateModels() log2.FunctionWithNoArgs {
	return log2.LogWraps(AutoMigrateModelsWithoutLogs)
}

func StartHTTPServer(router *gin.Engine) *http.Server {
	log.Printf("Started StartHTTPServer on %s", httpPort)
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

	log.Printf("Finished StartHTTPServer on %s", httpPort)
	return server
}

func GracefulShutdown(httpServer *http.Server) {
	log.Printf("Started GracefulShutdown")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}
	log.Println("HTTP server stopped")
	log.Printf("Started GracefulShutdown on %s")
}

func LoadEnvWithoutLogs() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func InitDatabaseWithoutLogs() {
	config.InitDB()
}

func AutoMigrateModelsWithoutLogs() {
	config.DB.AutoMigrate(&models.Wallet{}, &models.Transaction{})
}

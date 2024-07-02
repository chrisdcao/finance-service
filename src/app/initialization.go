package app

import (
	"context"
	"finance-service/config"
	"finance-service/models"
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

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func InitDatabase() {
	config.InitDB()
}

func AutoMigrateModels() {
	config.DB.AutoMigrate(&models.Wallet{}, &models.Transaction{})
}

func StartHTTPServer(router *gin.Engine) *http.Server {
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

func GracefulShutdown(httpServer *http.Server) {
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

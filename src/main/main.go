package main

import (
	"context"
	"finance-service/config"
	"finance-service/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

const (
	// TODO: Refactor this to .env file
	httpPort = ":3000" // Port for the HTTP server
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

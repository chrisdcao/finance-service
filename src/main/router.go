package main

import (
	"finance-service/controllers"
	"net/http"
)

func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Initialize controllers
	transactionController := controllers.NewTransactionController()
	walletController := controllers.NewWalletController()

	// Define transaction routes
	mux.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			transactionController.GetTransactions(w, r)
		} else if r.Method == http.MethodPost {
			transactionController.CreateTransaction(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Define wallet routes
	mux.HandleFunc("/wallets/update_balance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			walletController.UpdateBalance(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/wallets/convert_balance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			walletController.ConvertBalance(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

package main

import (
	"finance-service/controllers"
	"net/http"
)

// TODO: Move the constructed service to this place and have router using them?
func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Initialize controllers
	transactionController := controllers.NewEndUserControllerController()
	walletController := controllers.NewAdminController()

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
			walletController.Topup(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/wallets/convert_balance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			walletController.WalletTransfer(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

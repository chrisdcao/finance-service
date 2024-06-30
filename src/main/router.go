package main

import (
	"net/http"
)

// TODO: Move the constructed service to this place and have router using them?
func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Setup DI container
	container := NewContainer()

	// Define transaction routes
	// TODO: based on the standard of REST APIs (using same endpoint with diff methods for diff interactions)
	mux.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			container.EndUserController.GetTransactions(w, r)
			// TODO: define other endpoints + rbac here for other biz logics
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Define wallet routes
	mux.HandleFunc("/wallets/update_balance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			container.AdminController.Topup(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/wallets/convert_balance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			container.AdminController.WalletTransfer(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

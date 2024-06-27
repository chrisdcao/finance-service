package controllers

import (
	"encoding/json"
	"finance-service/models"
	"finance-service/services/wallet/transaction"
	"finance-service/utils"
	"net/http"
)

type EndUserController struct {
	TransactionService *transaction.TransactionWriteService
}

func NewEndUserControllerController() *EndUserController {
	return &EndUserController{TransactionService: transaction.NewTransactionWriteService()}
}

func (this *EndUserController) GetTransactionsByExternalWalletId(w http.ResponseWriter, r *http.Request) {
	transactions, err := this.TransactionService.GetTransactions(r.URL.Query().Get("uuid"))
	if err != nil {
		utils.Logger().Println("Error getting transactions:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		utils.Logger().Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (this *EndUserController) GetTransactionsByUserId(w http.ResponseWriter, r *http.Request) {
	transactions, err := this.TransactionService.GetTransactions(r.URL.Query().Get("uuid"))
	if err != nil {
		utils.Logger().Println("Error getting transactions:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		utils.Logger().Println("Error encoding response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (this *EndUserController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		utils.Logger().Println("Error decoding transaction:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := this.TransactionService.CreateTransaction(transaction); err != nil {
		utils.Logger().Println("Error creating transaction:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

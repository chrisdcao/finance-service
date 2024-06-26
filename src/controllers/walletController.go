package controllers

import (
	"context"
	"encoding/json"
	"finance-service/controllers/dto"
	walletservices "finance-service/services/wallet"
	exchangeservices "finance-service/services/wallet/balance/exchange"
	"finance-service/services/wallet/transaction"
	"finance-service/utils"
	"net/http"
)

type WalletController struct {
	WalletService      *walletservices.DefaultWalletWriteService
	ExchangeService    *exchangeservices.ExchangeService
	TransactionService *transaction.TransactionWriteService
}

func NewWalletController(walletService *walletservices.DefaultWalletWriteService, exchangeService *exchangeservices.ExchangeService) *WalletController {
	return &WalletController{WalletService: walletService, ExchangeService: exchangeService}
}

func (this *WalletController) TopUp(w http.ResponseWriter, r *http.Request) {
	var req dto.UpdateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if walletDto, err := this.WalletService.UpdateBalance(req); err != nil {
		utils.Logger().Println("Error updating balance:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Define response to return the updated wallet dto
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Balance updated successfully"})
}

func (this *WalletController) ConvertBalance(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UUID   string  `json:"uuid"`
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger().Println("Error decoding request:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	wallet, err := this.ExchangeService.ConvertBalance(ctx, req.UUID, req.Amount)
	if err != nil {
		utils.Logger().Println("Error converting balance:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

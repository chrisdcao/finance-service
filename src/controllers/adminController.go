package controllers

import (
	"encoding/json"
	"finance-service/controllers/dto/response"
	walletservices "finance-service/services/wallet"
	"finance-service/services/wallet/dto"
	"finance-service/utils"
	"net/http"
)

type AdminController struct {
	WalletService *walletservices.DefaultWalletWriteService
}

func NewAdminController(walletService *walletservices.DefaultWalletWriteService) *AdminController {
	return &AdminController{WalletService: walletService}
}

func (this *AdminController) Topup(w http.ResponseWriter, r *http.Request) {
	var req dto.WalletUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSONResponse(w, http.StatusBadRequest, err.Error(), nil, "Invalid request payload")
		return
	}

	walletDto, err := this.WalletService.UpdateBalance(nil, req)
	if err != nil {
		utils.Logger().Println("Error updating balance:", err)
		// TODO: Do we need a custom error code here?
		response.WriteJSONResponse(w, http.StatusInternalServerError, err.Error(), nil, "Error updating balance")
		return
	}

	response.WriteJSONResponse(w, http.StatusOK, "", walletDto, "Balance updated successfully")
}

func (this *AdminController) WalletTransfer(w http.ResponseWriter, r *http.Request) {
	var req dto.WalletTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteJSONResponse(w, http.StatusBadRequest, err.Error(), nil, "Invalid request payload")
		return
	}

	walletDto, err := this.WalletService.WalletTransfer(req)
	if err != nil {
		utils.Logger().Println("Error converting balance:", err)
		response.WriteJSONResponse(w, http.StatusInternalServerError, err.Error(), nil, "Error updating balance")
		return
	}

	response.WriteJSONResponse(w, http.StatusOK, "", walletDto, "Balance updated successfully")
}

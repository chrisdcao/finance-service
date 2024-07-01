package controllers

import (
	"encoding/json"
	"finance-service/controllers/dto/response"
	walletservices "finance-service/services/wallet"
	"finance-service/services/wallet/dto"
	"finance-service/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminController struct {
	WalletService *walletservices.DefaultWalletWriteService
}

func NewAdminController(walletService *walletservices.DefaultWalletWriteService) *AdminController {
	return &AdminController{WalletService: walletService}
}

// Topup godoc
// @Summary Top up wallet balance
// @Description Top up the balance of a wallet
// @Tags wallets
// @Accept json
// @Produce json
// @Param topupRequest body dto.WalletUpdateRequest true "Topup Request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /wallets/update_balance [post]
func (this *AdminController) Topup(ctx *gin.Context) {
	var req dto.WalletUpdateRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		response.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid request payload")
		return
	}

	walletDto, err := this.WalletService.UpdateBalance(nil, req)
	if err != nil {
		utils.Logger().Println("Error updating balance:", err)
		// TODO: Do we need a custom error code here?
		response.WriteJSONResponse(ctx, http.StatusInternalServerError, err.Error(), nil, "Error updating balance")
		return
	}

	response.WriteJSONResponse(ctx, http.StatusOK, "", walletDto, "Balance updated successfully")
}

// WalletTransfer godoc
// @Summary Convert wallet balance
// @Description Transfer balance from one wallet to another
// @Tags wallets
// @Accept json
// @Produce json
// @Param convertRequest body dto.WalletTransferRequest true "Convert Request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /wallets/convert_balance [post]
func (this *AdminController) WalletTransfer(ctx *gin.Context) {
	var req dto.WalletTransferRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		response.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid request payload")
		return
	}

	postTransferredWallets, err := this.WalletService.WalletTransfer(req)
	if err != nil {
		utils.Logger().Println("Error converting balance:", err)
		response.WriteJSONResponse(ctx, http.StatusInternalServerError, err.Error(), nil, "Error updating balance")
		return
	}

	response.WriteJSONResponse(ctx, http.StatusOK, "", postTransferredWallets, "Balance updated successfully")
}

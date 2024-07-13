package controllers

import (
	"encoding/json"
	log2 "finance-service/common/log"
	logDto "finance-service/common/log/dto"
	"finance-service/common/web"
	"finance-service/controllers/wallet/dto/request"
	response2 "finance-service/controllers/wallet/dto/response"
	walletservices "finance-service/services/wallet"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AdminController struct {
	WalletService *walletservices.DefaultWalletService
	Logger        *log2.CommonLogger
}

func NewAdminController(walletService *walletservices.DefaultWalletService) *AdminController {
	return &AdminController{
		WalletService: walletService,
		Logger:        log2.NewCommonLogger(),
	}
}

// Topup godoc
// @Summary Topup Wallet
// @Description Topup a wallet with the specified amount
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param   request body request.WalletUpdateRequest true "Topup request payload"
// @Success 200 {object} response2.WalletUpdateResponse
// @Failure 400 {object} web.Response "Invalid request payload"
// @Failure 500 {object} web.Response "Error updating balance"
// @Router /admin/topup [post]
func (this *AdminController) Topup(ctx *gin.Context) {
	var req request.WalletUpdateRequest

	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		web.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid request payload")
	}

	transactionDto, err := this.WalletService.UpdateBalance(ctx, nil, req)

	if err != nil {
		this.Logger.Log(logDto.LogEntry{
			Level:   logrus.ErrorLevel,
			TraceID: log2.GetTraceIDFromGinContextOrUnknown(ctx),
			Event:   "process_request",
			Message: "Request processing failed",
			Context: map[string]interface{}{
				"error": errors.WithStack(err),
			},
		})

		// TODO: Do we need a custom error code here? Discuss with Duc Huy
		web.WriteJSONResponse(ctx, http.StatusInternalServerError, errors.WithStack(err).Error(), nil, "Error updating balance")
	}

	resp := response2.NewWalletUpdateResponse(*transactionDto)

	web.WriteJSONResponse(ctx, http.StatusOK, "", resp, "Balance updated successfully")
}

// WalletTransfer godoc
// @Summary Transfer between wallets
// @Description Transfer funds between wallets
// @Tags Admin
// @Accept  json
// @Produce  json
// @Param   request body request.WalletTransferRequest true "Wallet transfer request payload"
// @Success 200 {object} response2.WalletTransferResponse
// @Failure 400 {object} web.Response "Invalid request payload"
// @Failure 500 {object} web.Response "Error processing transfer"
// @Router /admin/transfer [post]
func (this *AdminController) WalletTransfer(ctx *gin.Context) {
	var req request.WalletTransferRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		web.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid request payload")
	}

	createdTransactions, err := this.WalletService.WalletTransfer(ctx, req)
	if err != nil {
		this.Logger.Log(logDto.LogEntry{
			Level:   logrus.ErrorLevel,
			TraceID: log2.GetTraceIDFromGinContextOrUnknown(ctx),
			Event:   "process_request",
			Message: "Request processing failed",
			Context: map[string]interface{}{
				"error": errors.WithStack(err),
			},
		})

		web.WriteJSONResponse(ctx, http.StatusInternalServerError, errors.WithStack(err).Error(), nil, "Error updating balance")
	}

	resp := response2.NewWalletTransferResponse(createdTransactions)

	web.WriteJSONResponse(ctx, http.StatusOK, "", resp, "Balance updated successfully")
}

package controllers

import (
	"encoding/json"
	"finance-service/controllers/wallet/dto/request"
	response2 "finance-service/controllers/wallet/dto/response"
	walletservices "finance-service/services/wallet"
	"finance-service/utils/log"
	logDto "finance-service/utils/log/dto"
	"finance-service/utils/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AdminController struct {
	WalletService *walletservices.DefaultWalletService
	Logger        *log.CommonLogger
}

func NewAdminController(walletService *walletservices.DefaultWalletService) *AdminController {
	return &AdminController{
		WalletService: walletService,
		Logger:        log.NewCommonLogger(),
	}
}

// TODO: The update request should be receiving an user_id and return the transaction that was made?
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
	var req request.WalletUpdateRequest

	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		web.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid request payload")
	}

	transactionDto, err := this.WalletService.UpdateBalance(ctx, nil, req)

	if err != nil {
		this.Logger.Log(logDto.LogEntry{
			Level:   logrus.ErrorLevel,
			TraceID: log.GetTraceIDFromGinContextOrUnknown(ctx),
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

// This should be returning 2 transaction within a transaction array?
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
	var req request.WalletTransferRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		web.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid request payload")
	}

	createdTransactions, err := this.WalletService.WalletTransfer(ctx, req)
	if err != nil {
		this.Logger.Log(logDto.LogEntry{
			Level:   logrus.ErrorLevel,
			TraceID: log.GetTraceIDFromGinContextOrUnknown(ctx),
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

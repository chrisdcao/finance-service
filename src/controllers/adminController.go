package controllers

import (
	"encoding/json"
	request2 "finance-service/controllers/dto/request"
	"finance-service/controllers/dto/response"
	walletservices "finance-service/services/wallet"
	"finance-service/utils/log"
	logDto "finance-service/utils/log/dto"
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
	var req request2.WalletUpdateRequest

	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		response.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid request payload")
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
		response.WriteJSONResponse(ctx, http.StatusInternalServerError, errors.WithStack(err).Error(), nil, "Error updating balance")
	}

	resp := response.NewWalletUpdateResponse(*transactionDto)

	response.WriteJSONResponse(ctx, http.StatusOK, "", resp, "Balance updated successfully")
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
	var req request2.WalletTransferRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
		response.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid request payload")
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

		response.WriteJSONResponse(ctx, http.StatusInternalServerError, errors.WithStack(err).Error(), nil, "Error updating balance")
	}

	resp := response.NewWalletTransferResponse(createdTransactions)

	response.WriteJSONResponse(ctx, http.StatusOK, "", resp, "Balance updated successfully")
}

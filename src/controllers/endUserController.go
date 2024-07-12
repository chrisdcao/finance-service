package controllers

import (
	"finance-service/controllers/dto/response"
	"finance-service/services/transaction"
	transactiondto "finance-service/services/transaction/dto"
	"finance-service/utils/log"
	"finance-service/utils/log/dto"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

// EndUserController handles end-user related operations
type EndUserController struct {
	TransactionReadService *transaction.TransactionReadService
	Logger                 *log.CommonLogger
}

// NewEndUserController creates a new EndUserController
func NewEndUserController(service *transaction.TransactionReadService) *EndUserController {
	return &EndUserController{
		TransactionReadService: service,
		Logger:                 log.NewCommonLogger(),
	}
}

// GetTransactions godoc
// @Summary Get transactions
// @Description Get all transactions based on filters
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param walletType query string false "Wallet TransactionType"
// @Param actionType query string false "Action TransactionType"
// @Param amount query string false "DiffAmount"
// @Param fromTime query int64 false "From Time (ms)"
// @Param toTime query int64 false "To Time (ms)"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /transactions/get [get]
func (this *EndUserController) GetTransactions(ctx *gin.Context) {
	traceId := log.ExtractTraceIDFromContextOrUnknown(ctx)

	var params transactiondto.GetTransactionsRequest

	if err := ctx.ShouldBindQuery(&params); err != nil {
		this.Logger.Log(dto.LogEntry{
			Level:      logrus.ErrorLevel,
			Event:      "parsing_query_params",
			Message:    "Error parsing query params",
			Context:    map[string]interface{}{"params": params},
			StackTrace: errors.WithStack(err),
		})
		response.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid query parameters")
		return
	}

	transactions, err := this.TransactionReadService.GetTransactions(&params)
	if err != nil {
		this.Logger.Log(dto.LogEntry{
			Level:   logrus.ErrorLevel,
			TraceID: traceId,
			Event:   "getting_transactions",
			Message: "Error getting transactions",
			Context: map[string]interface{}{
				"params": params,
			},
			StackTrace: errors.WithStack(err),
		})
		response.WriteJSONResponse(ctx, http.StatusInternalServerError, err.Error(), nil, "Error retrieving transactions")
		return
	}

	response.WriteJSONResponse(ctx, http.StatusOK, "", transactions, "Transactions retrieved successfully")
}

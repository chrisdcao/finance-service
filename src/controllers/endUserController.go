package controllers

import (
	log2 "finance-service/common/log"
	"finance-service/common/log/dto"
	"finance-service/common/web"
	transactiondto "finance-service/controllers/transaction/dto/request"
	"finance-service/controllers/transaction/dto/response"
	"finance-service/services/transaction"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

// EndUserController handles end-user related operations
type EndUserController struct {
	TransactionReadService *transaction.TransactionReadService
	Logger                 *log2.CommonLogger
}

// NewEndUserController creates a new EndUserController
func NewEndUserController(service *transaction.TransactionReadService) *EndUserController {
	return &EndUserController{
		TransactionReadService: service,
		Logger:                 log2.NewCommonLogger(),
	}
}

// GetTransactions godoc
// @Summary Get Transactions
// @Description Retrieve transactions based on query parameters
// @Tags EndUser
// @Accept  json
// @Produce  json
// @Param   wallet_type query string false "Wallet Type"
// @Param   action_type query string false "Action Type"
// @Param   amount query number false "Amount"
// @Param   from_time query string false "From Time"
// @Param   to_time query string false "To Time"
// @Param   uuid query string false "UUID"
// @Success 200 {object} response.GetTransactionsResponse
// @Failure 400 {object} web.Response "Invalid query parameters"
// @Failure 500 {object} web.Response "Error retrieving transactionDtos"
// @Router /enduser/transactions [get]
func (this *EndUserController) GetTransactions(ctx *gin.Context) {
	traceId := log2.ExtractTraceIDFromContextOrUnknown(ctx)

	var params transactiondto.GetTransactionsRequest

	if err := ctx.ShouldBindQuery(&params); err != nil {
		this.Logger.Log(dto.LogEntry{
			Level:      logrus.ErrorLevel,
			Event:      "parsing_query_params",
			Message:    "Error parsing query params",
			Context:    map[string]interface{}{"params": params},
			StackTrace: errors.WithStack(err),
		})
		web.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid query parameters")
		return
	}

	transactionDtos, err := this.TransactionReadService.GetTransactions(params)
	if err != nil {
		this.Logger.Log(dto.LogEntry{
			Level:   logrus.ErrorLevel,
			TraceID: traceId,
			Event:   "getting_transactions",
			Message: "Error getting transactionDtos",
			Context: map[string]interface{}{
				"params": params,
			},
			StackTrace: errors.WithStack(err),
		})
		web.WriteJSONResponse(ctx, http.StatusInternalServerError, err.Error(), nil, "Error retrieving transactionDtos")
		return
	}

	web.WriteJSONResponse(ctx, http.StatusOK, "", response.NewGetTransactionsResponse(transactionDtos), "Transactions retrieved successfully")
}

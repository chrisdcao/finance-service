package controllers

import (
	"finance-service/controllers/dto/response"
	"finance-service/services/wallet/transaction"
	transactiondto "finance-service/services/wallet/transaction/dto"
	"finance-service/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// EndUserController handles end-user related operations
type EndUserController struct {
	TransactionReadService *transaction.TransactionReadService
}

// NewEndUserController creates a new EndUserController
func NewEndUserController(service *transaction.TransactionReadService) *EndUserController {
	return &EndUserController{TransactionReadService: service}
}

// GetTransactions godoc
// @Summary Get transactions
// @Description Get all transactions based on filters
// @Tags transactions
// @Accept  json
// @Produce  json
// @Param walletType query string false "Wallet TransactionType"
// @Param actionType query string false "Action TransactionType"
// @Param amount query string false "Amount"
// @Param fromTime query int64 false "From Time (ms)"
// @Param toTime query int64 false "To Time (ms)"
// @Success 200 {object} dto.GenericResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v2/user/transaction [get]
func (c *EndUserController) GetTransactions(ctx *gin.Context) {
	var params transactiondto.GetTransactionsRequest
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.Logger().Println("Error parsing query params:", err)
		response.WriteJSONResponse(ctx, http.StatusBadRequest, err.Error(), nil, "Invalid query parameters")
		return
	}

	transactions, err := c.TransactionReadService.GetTransactions(&params)
	if err != nil {
		utils.Logger().Println("Error getting transactions:", err)
		response.WriteJSONResponse(ctx, http.StatusInternalServerError, err.Error(), nil, "Error retrieving transactions")
		return
	}

	response.WriteJSONResponse(ctx, http.StatusOK, "", transactions, "Transactions retrieved successfully")
}

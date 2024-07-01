package controllers

import (
	"finance-service/controllers/dto/response"
	"finance-service/services/wallet/transaction"
	"finance-service/services/wallet/transaction/dto"
	"finance-service/utils"
	"net/http"
)

type EndUserController struct {
	TransactionReadService *transaction.TransactionReadService
}

func NewEndUserControllerController(service *transaction.TransactionReadService) *EndUserController {
	return &EndUserController{TransactionReadService: service}
}

// TODO: Add here possible query params and path variablAes
// Like old time REST Java Spring
// Based on that we return the transaction (w/ or w/o filter, filters are the params passed in)
// we always have user on jwt -> we could be :
// /api/v2/user/transaction
// headers: jwt
// params: ?walletType=&actionType=&amount=&fromTime=1231245125&toTime=... (time will be in ms, utc0)
func (this *EndUserController) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var params dto.GetTransactionsRequest
	err := utils.ParseQueryParams(r, params)
	if err != nil {
		utils.Logger().Println("Error parsing query params:", err)
		response.WriteJSONResponse(w, http.StatusBadRequest, err.Error(), nil, "Invalid query parameters")
		return
	}

	transactions, err := this.TransactionReadService.GetTransactions(&params)
	if err != nil {
		utils.Logger().Println("Error getting transactions:", err)
		response.WriteJSONResponse(w, http.StatusInternalServerError, err.Error(), nil, "Error retrieving transactions")
		return
	}

	response.WriteJSONResponse(w, http.StatusOK, "", transactions, "Transactions retrieved successfully")
}

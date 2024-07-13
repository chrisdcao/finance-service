package request

type WalletTransferRequest struct {
	UserId string  `json:"userId"`
	Amount float64 `json:"amount"`
}

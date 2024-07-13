package request

import "time"

// Request struct for transaction query parameters
type GetTransactionsRequest struct {
	WalletType string    `query:"walletType"`
	ActionType string    `query:"actionType"`
	Amount     float64   `query:"amount"`
	FromTime   time.Time `query:"fromTime"`
	ToTime     time.Time `query:"toTime"`
	UUID       string    `query:"uuid"`
}

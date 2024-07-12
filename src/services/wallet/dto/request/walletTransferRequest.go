package request

type WalletTransferRequest struct {
	ExternalFromWalletId string  `json:"externalFromWalletId"`
	ExternalToWalletId   string  `json:"externalToWalletId"`
	Amount               float64 `json:"amount"`
}

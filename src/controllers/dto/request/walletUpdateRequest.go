package request

type WalletUpdateRequest struct {
	ExternalWalletId string  `json:"externalWalletId"`
	Amount           float64 `json:"amount"`
	UpdateType       string  `json:"updateType"`
	Content          string  `json:"content"`
	WalletType       string  `json:"walletType"`
}

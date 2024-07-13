package request

type WalletUpdateRequest struct {
	// TODO: Check with Huy data type of userId
	UserId     string  `json:"userId"`
	Amount     float64 `json:"amount"`
	UpdateType string  `json:"updateType"`
	Content    string  `json:"content"`
	WalletType string  `json:"walletType"`
}

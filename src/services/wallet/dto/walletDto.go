package dto

type WalletDto struct {
	ID         uint    `gorm:"primaryKey"`
	UserId     int     `json:"user_id"`
	Balance    float64 `json:"balance"`
	WalletType string  `json:"wallet_type"`
}

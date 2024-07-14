package enums

import "github.com/pkg/errors"

// WalletType represents different types of wallets
type WalletType int

const (
	ASMWallet WalletType = iota
	VNDWallet
)

// Mappings for WalletType and WalletType to their string representations
var walletTypeToString = map[WalletType]string{
	ASMWallet: "ASM Wallet",
	VNDWallet: "VND Wallet",
}

var stringToWalletType = map[string]WalletType{
	"ASM Wallet": ASMWallet,
	"VND Wallet": VNDWallet,
}

func (this WalletType) String() string {
	return walletTypeToString[this]
}

func ParseStringToWalletType(str string) (WalletType, error) {
	t, ok := stringToWalletType[str]
	if !ok {
		return 0, errors.New("invalid wallet type: " + str)
	}
	return t, nil
}

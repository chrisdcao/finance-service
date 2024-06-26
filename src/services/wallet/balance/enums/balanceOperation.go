package enums

import "fmt"

type BalanceOperation int

const (
	ASMWalletDebit BalanceOperation = iota
	ASMWalletTopup
	VNDWalletDebit
	VNDWalletTopup
)

var topupTypeToString = map[BalanceOperation]string{
	ASMWalletDebit: "ASM Wallet Debit",
	ASMWalletTopup: "ASM Wallet Topup",
	VNDWalletDebit: "VND Wallet Debit",
	VNDWalletTopup: "VND Wallet Topup",
}

var stringToTopupType = map[string]BalanceOperation{
	"ASM Wallet Debit": ASMWalletDebit,
	"ASM Wallet Topup": ASMWalletTopup,
	"VND Wallet Debit": VNDWalletDebit,
	"VND Wallet Topup": VNDWalletTopup,
}

func (this BalanceOperation) String() string {
	return topupTypeToString[this]
}

func ParseTopupType(str string) (BalanceOperation, error) {
	t, ok := stringToTopupType[str]
	if !ok {
		return 0, fmt.Errorf("invalid topup type: %s", str)
	}
	return t, nil
}

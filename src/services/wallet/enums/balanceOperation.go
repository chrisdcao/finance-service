package enums

import "github.com/pkg/errors"

// BalanceOperation represents different balance operations
type BalanceOperation int

const (
	Credit BalanceOperation = iota
	Debit
)

var balanceOperationToString = map[BalanceOperation]string{
	Credit: "Credit",
	Debit:  "Debit",
}

var stringToBalanceOperation = map[string]BalanceOperation{
	"Credit": Credit,
	"Debit":  Debit,
}

func (this BalanceOperation) String() string {
	return balanceOperationToString[this]
}

func ParseStringToBalanceOperation(str string) (BalanceOperation, error) {
	t, ok := stringToBalanceOperation[str]
	if !ok {
		return 0, errors.New("invalid balance operation: " + str)
	}
	return t, nil
}

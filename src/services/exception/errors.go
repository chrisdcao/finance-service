package exception

import "errors"

var (
	ErrTransactionNotFound = errors.New("transactions not found")
)

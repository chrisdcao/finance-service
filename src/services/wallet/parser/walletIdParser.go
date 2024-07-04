package parser

import (
	"context"
	"finance-service/services/cryptography"
	"github.com/pkg/errors"
	"strconv"
)

var DUMMY_VALUE = uint(0)

type WalletIdParser struct {
}

func (this *WalletIdParser) ParseFromEncryption(ctx context.Context, encryptedWalletId string) (uint, error) {
	walletIdStr, err := cryptography.Decrypt(encryptedWalletId)
	if err != nil {
		return DUMMY_VALUE, errors.Wrap(err, "failed to decrypt wallet id: "+encryptedWalletId)
	}

	walletIdLong, err := strconv.ParseUint(walletIdStr, 10, 64)
	if err != nil {
		return DUMMY_VALUE, errors.Wrap(err, "failed to parse wallet id to int from string: "+walletIdStr)
	}

	return uint(walletIdLong), nil
}

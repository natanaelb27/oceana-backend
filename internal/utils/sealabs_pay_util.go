package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"git.garena.com/azmi-adhani/group-project-backend/internal/constants"
	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

func ValidateSealabsPayAccountNumber(accountNumber string) bool {
	return len(accountNumber) == constants.SealabsPayAccountNumberLength
}

func HashPaymentInfoSHA256HMAC(cardNumber, amount, merchantCode, key string) (*string, error) {
	hasher := hmac.New(sha256.New, []byte(key))
	signatureStringToEncrypt := fmt.Sprintf("%s:%s:%s", cardNumber, amount, merchantCode)
	hasher.Write([]byte(signatureStringToEncrypt))
	signatureBytes := hasher.Sum(nil)
	signatureString := hex.EncodeToString(signatureBytes)
	if len(signatureString) == 0 {
		return nil, e.ErrHashingFailed
	}

	return &signatureString, nil
}

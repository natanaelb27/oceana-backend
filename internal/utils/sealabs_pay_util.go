package utils

import (
	"strconv"
)

func ValidateSealabsPayAccountNumber(accountNumber string) bool {
	if len(accountNumber) != 16 {
		return false
	}
	if _, err := strconv.ParseInt(accountNumber, 10, 64); err != nil {
		return false
	}
	return true
}

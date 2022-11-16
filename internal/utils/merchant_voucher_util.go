package utils

import "unicode"

func ValidateMerchantVoucherCode(voucherCode string) bool {
	if len(voucherCode) > 5 {
		return false
	}

	for _, char := range voucherCode {
		switch {
		case !unicode.IsNumber(char):
			if unicode.IsLower(char) {
				return false
			}
		}
	}

	return true
}

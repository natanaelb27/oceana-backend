package utils

import (
	"github.com/natanaelb27/oceana-backend/internal/constants"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

func ValidatePmpDiscountAndType(discountType string, discount uint) error {
	if discountType != constants.Nominal && discountType != constants.Percentage {
		return e.ErrInvalidDiscountType
	}
	if discount <= 0 {
		return e.ErrInvalidDiscount
	}
	if discountType == constants.Percentage {
		if discount > 100 {
			return e.ErrExceedDiscountPercentage
		}
	}
	return nil
}

func ValidatePmpQuota(currentQuota uint, maximumQuota, maximumUserLimit int) error {
	if maximumQuota == -1 {
		return nil
	}
	if maximumQuota < int(currentQuota) {
		return e.ErrMaximumQuotaLowerThanCurrentQuota
	}
	if maximumQuota < maximumUserLimit {
		return e.ErrMaximumQuotaLowerThanMaximumUserLimit
	}
	return nil
}

package utils

import (
	"time"

	e "git.garena.com/azmi-adhani/group-project-backend/pkg/errors"
)

func ValidateMerchantPromotionDatetime(startDatetime, endDatetime time.Time) error {
	if startDatetime.Before(time.Now()) {
		return e.ErrPromotionStartDatetimeEarlierThanNow
	}

	if endDatetime.Before(startDatetime) {
		return e.ErrPromotionEndDatetimeEarlierThanVoucherStartDatetime
	}

	if int(time.Until(startDatetime).Hours()/24) <= 0 {
		return e.ErrPromotionStartDatetimeMinimumOneDayFromNow
	}

	if int(endDatetime.Sub(startDatetime).Hours()) <= 0 {
		return e.ErrPromotionPeriodMinimumOneHour
	}

	if int(endDatetime.Sub(startDatetime).Hours()/24) > 180 {
		return e.ErrPromotionPeriodExceedMaximumDays
	}

	return nil
}

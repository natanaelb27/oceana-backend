package utils

import (
	"time"

	"github.com/natanaelb27/oceana-backend/internal/constants"
	e "github.com/natanaelb27/oceana-backend/pkg/errors"
)

func ValidateMpCreateDatetime(startDatetime, endDatetime time.Time) error {
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

	if int(endDatetime.Sub(startDatetime).Hours()/24) > 90 {
		return e.ErrPromotionPeriodExceedMaximumDays
	}

	return nil
}

func ValidateMpUpdateDatetime(startDatetime, endDatetime, currentStartDatetime, currentEndDatetime time.Time, promotionStatus string) error {
	// promotion ended, no updates allowed
	if promotionStatus == constants.EndedPromotionOrVoucher {
		return e.ErrPromotionHasEnded
	}
	// promotion ongoing, cannot update start time
	if promotionStatus == constants.OngoingPromotionOrVoucher {
		if startDatetime != currentStartDatetime {
			return e.ErrOngoingStartDatetimeUpdated
		}
	}
	// promotion upcoming & ongoing must be inbetween current start/end datetime
	if startDatetime.Before(currentStartDatetime) || endDatetime.After(currentEndDatetime) {
		return e.ErrPromotionDatetimeUpdatedBeyondRange
	}

	return nil
}

func ValidateMpProductIDList(mpRequestProductIDList, pmpRequestProductIDList []int) error {
	if len(mpRequestProductIDList) != len(pmpRequestProductIDList) {
		return e.ErrPromotionUpdateProductIDListMismatch
	}
	for i := 0; i < len(mpRequestProductIDList); i++ {
		if mpRequestProductIDList[i] != pmpRequestProductIDList[i] {
			return e.ErrPromotionUpdateProductIDListMismatch
		}
	}
	return nil
}

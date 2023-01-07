package utils

import (
	"time"

	"github.com/natanaelb27/oceana-backend/internal/constants"
)

func IsDateMoreThanOneDay(date time.Time) bool {
	now := time.Now()
	newDate := date.Add(24 * time.Hour)

	return now.Before(newDate)
}

func PromotionOrVoucherStatus(startDatetime, endDatetime time.Time) string {
	now := time.Now()
	if now.Before(startDatetime) {
		return constants.UpcomingPromotionOrVoucher
	} else if now.After(endDatetime) {
		return constants.EndedPromotionOrVoucher
	} else {
		return constants.OngoingPromotionOrVoucher
	}
}

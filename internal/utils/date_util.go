package utils

import "time"

func IsDateMoreThanOneDay(date time.Time) bool {
	now := time.Now()
	newDate := date.Add(24 * time.Hour)

	return now.Before(newDate)
}

func PromotionOrVoucherStatus(startDatetime, endDatetime time.Time) string {
	now := time.Now()
	if now.Before(startDatetime) {
		return "upcoming"
	} else if now.After(endDatetime) {
		return "ended"
	} else {
		return "ongoing"
	}
}

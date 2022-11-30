package utils

import "time"

func IsDateMoreThanOneDay(date time.Time) bool {
	now := time.Now()
	newDate := date.Add(24 * time.Hour)

	return now.Before(newDate)
}

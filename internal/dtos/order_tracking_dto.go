package dtos

import "time"

type OrderTrackingHistory struct {
	Datetime    time.Time `json:"datetime"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

package dtos

import (
	"time"

	"github.com/google/uuid"
)

type OrderTrackingHistory struct {
	Datetime    time.Time `json:"datetime"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type COTUOrderTracking struct {
	ID              uint
	HistoryDatetime time.Time
	Title           string
	Description     string
	OrderID         uuid.UUID
	OrderStatus     string
}

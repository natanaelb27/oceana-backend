package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderTracking struct {
	ID              uint
	OrderID         uuid.UUID
	Order           Order
	HistoryDatetime time.Time
	Title           string
	Description     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}

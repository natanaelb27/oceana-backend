package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDetail struct {
	ID                      uint
	OrderID                 uuid.UUID
	Order                   Order
	ProductSnapshotID       uuid.UUID
	ProductSnapshot         ProductSnapshot
	Quantity                uint
	TotalPrice              uint
	TotalPriceAfterDiscount uint
	Status                  string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               gorm.DeletedAt
}

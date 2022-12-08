package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDetail struct {
	ID                      uuid.UUID
	OrderID                 uuid.UUID
	Order                   Order
	Message                 string
	ProductSnapshotID       uint
	ProductSnapshot         ProductSnapshot
	Quantity                uint
	TotalPrice              uint
	TotalPriceAfterDiscount float64
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               gorm.DeletedAt
}

// func (m *OrderDetail) BeforeCreate(tx *gorm.DB) (err error) {
// 	m.ID = uuid.New()
// 	return
// }

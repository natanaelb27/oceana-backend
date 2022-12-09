package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID                        uuid.UUID
	UserID                    uuid.UUID
	User                      User
	TransactionDatetime       time.Time
	TotalPrice                float64
	TotalPriceAfterDiscount   float64
	TotalPriceWithShippingFee float64
	PaymentMethodID           uint
	PaymentMethod             PaymentMethod
	IsPaid                    bool
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	DeletedAt                 gorm.DeletedAt
}

// func (m *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
// 	m.ID = uuid.New()
// 	return
// }

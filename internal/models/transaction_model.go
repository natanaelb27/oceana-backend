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
	TotalPrice                uint
	TotalPriceAfterDiscount   uint
	TotalPriceWithShippingFee uint
	PaymentMethodID           uint
	PaymentMethod             PaymentMethod
	AddressSnapshotID         uint
	AddressSnapshot           AddressSnapshot
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	DeletedAt                 gorm.DeletedAt
}

func (m *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

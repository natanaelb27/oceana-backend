package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID                        uuid.UUID
	TxID                      string
	UserID                    uuid.UUID
	User                      User
	TransactionDate           time.Time
	TotalPrice                int
	TotalPriceAfterDiscount   int
	TotalPriceWithShippingFee int
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

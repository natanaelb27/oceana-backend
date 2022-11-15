package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID                      uuid.UUID
	TransactionID           uuid.UUID
	Transaction             Transaction
	MerchantID              uint
	Merchant                Merchant
	CourierServiceID        uint
	CourierService          CourierService
	TotalPrice              int
	TotalPriceAfterDiscount int
	ShippingFee             int
	Status                  string
	ClaimStatus             bool
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               gorm.DeletedAt
}

func (m *Order) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

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
	AddressSnapshotID       uint
	AddressSnapshot         AddressSnapshot
	CourierServiceID        uint
	CourierService          CourierService
	TotalPrice              uint
	TotalPriceAfterDiscount uint
	ShippingFee             uint
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

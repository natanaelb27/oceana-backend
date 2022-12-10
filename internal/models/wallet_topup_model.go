package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletTopUp struct {
	ID              uuid.UUID
	WalletID        uuid.UUID
	Wallet          Wallet
	SealabsPayID    uuid.UUID
	SealabsPay      SealabsPay
	TxID            string
	Amount          uint
	TransactionDate time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}

func (m *WalletTopUp) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

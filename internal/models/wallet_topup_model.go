package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletTopup struct {
	ID              uuid.UUID
	WalletID        uuid.UUID
	Wallet          Wallet
	SealabsPayID    uuid.UUID
	SealabsPay      SealabsPay
	TxID            string
	Amount          int
	TransactionDate time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt
}

func (m *WalletTopup) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

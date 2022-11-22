package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletTransactionPayment struct {
	ID            uint
	TransactionID uuid.UUID
	Transaction   Transaction
	WalletID      uuid.UUID
	Wallet        Wallet
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

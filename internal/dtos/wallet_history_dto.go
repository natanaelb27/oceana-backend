package dtos

import (
	"time"

	"github.com/google/uuid"
)

type WalletTransactionHistory struct {
	ID              uuid.UUID
	WalletID        uuid.UUID
	Amount          int
	TransactionDate time.Time
	TransactionType string
	Type            string
}

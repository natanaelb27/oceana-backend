package dtos

import (
	"time"

	"github.com/google/uuid"
)

type WalletTransactionHistory struct {
	ID              uuid.UUID `json:"id"`
	WalletID        uuid.UUID `json:"wallet_id"`
	Amount          int       `json:"amount"`
	TransactionDate time.Time `json:"transaction_date"`
	TransactionType string    `json:"transcation_type"`
	Type            string    `json:"type"`
}

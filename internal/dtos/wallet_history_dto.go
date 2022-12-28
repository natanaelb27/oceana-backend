package dtos

import (
	"time"

	"github.com/google/uuid"
)

type WalletTransactionHistory struct {
	ID                  uuid.UUID `json:"id"`
	WalletID            uuid.UUID `json:"wallet_id"`
	Amount              int       `json:"amount"`
	TransactionDatetime time.Time `json:"transaction_datetime"`
	TransactionType     string    `json:"transaction_type"`
	Type                string    `json:"type"`
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SealabsPayTransactionPayment struct {
	ID                      uint
	TransactionID           uuid.UUID
	Transaction             Transaction
	SealabsPayAccountNumber string
	TxID                    string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               gorm.DeletedAt
}

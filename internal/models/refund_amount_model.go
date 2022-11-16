package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefundAmount struct {
	ID                 uuid.UUID
	RefundID           uuid.UUID
	Refund             Refund
	AmountUserGet      int
	AmountMerchantLost int
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt
}

func (m *RefundAmount) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefundDetail struct {
	ID                 uuid.UUID
	RefundID           uuid.UUID
	Refund             Refund
	OrderDetailID      uuid.UUID
	OrderDetail        OrderDetail
	Quantity           uint
	AmountUserGet      float64
	AmountMerchantLost float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt
}

func (m *RefundDetail) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

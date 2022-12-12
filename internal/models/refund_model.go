package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Refund struct {
	ID                       uuid.UUID
	OrderID                  uuid.UUID
	Order                    Order
	Reason                   string
	RefundRequestDatetime    time.Time
	RefundExpirationDatetime time.Time
	Status                   string
	RejectedDatetime         time.Time
	CreatedAt                time.Time
	UpdatedAt                time.Time
	DeletedAt                gorm.DeletedAt
}

func (m *Refund) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

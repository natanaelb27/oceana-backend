package models

import (
	"time"

	"gorm.io/gorm"
)

type MerchantConfiguration struct {
	ID           uint
	MerchantID   uint `gorm:"uniqueIndex:idx_merchant_id;"`
	PreorderDay  uint
	IsPrintLabel bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

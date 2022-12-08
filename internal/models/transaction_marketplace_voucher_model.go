package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionMarketplaceVoucher struct {
	ID                   uint
	TransactionID        uuid.UUID
	Transaction          Transaction
	MarketplaceVoucherID uint
	MarketplaceVoucher   MarketplaceVoucher
	TotalDiscount        float64
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt
}

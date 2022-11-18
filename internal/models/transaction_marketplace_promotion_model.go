package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionMarketplacePromotion struct {
	ID                     uint
	TransactionID          uuid.UUID
	Transaction            Transaction
	MarketplacePromotionID uint
	MarketplacePromotion   MarketplacePromotion
	TotalDiscount          uint
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt
}

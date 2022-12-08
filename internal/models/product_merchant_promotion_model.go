package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductMerchantPromotion struct {
	ID                    uint
	MerchantPromotionId   uint
	MerchantPromotion     MerchantPromotion
	ProductID             uint
	Product               Product
	DiscountType          string
	Discount              uint
	MaximumDiscountAmount uint
	MaximumQuota          int
	MaximumUserLimit      int
	CurrentQuota          uint
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt
}

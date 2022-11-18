package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductMerchantPromotion struct {
	ID                  uint
	MerchantPromotionId uint
	MerchantPromotion   MerchantPromotion
	ProductID           uint
	Product             Product
	MaximumQuota        uint
	MaximumUserLimit    uint
	CurrentQuota        uint
	Status              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
}

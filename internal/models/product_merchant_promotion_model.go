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
	MaximumQuota        int
	MaximumUserLimit    int
	CurrentQuota        int
	Status              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
}

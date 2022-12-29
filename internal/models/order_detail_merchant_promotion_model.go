package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDetailMerchantPromotion struct {
	ID                  uint
	OrderDetailID       uuid.UUID
	OrderDetail         OrderDetail
	ProductMerchantPromotionID uint
	ProductMerchantPromotion   ProductMerchantPromotion
	TotalDiscount       float64
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
}

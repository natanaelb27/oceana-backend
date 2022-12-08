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
	MerchantPromotionID uint
	MerchantPromotion   MerchantPromotion
	TotalDiscount       float64
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
}

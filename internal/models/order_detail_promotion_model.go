package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderDetailPromotion struct {
	ID                  uint
	OrderDetailID       uuid.UUID
	OrderDetail         OrderDetail
	MerchantPromotionID uint
	MerchantPromotion   MerchantPromotion
	TotalDiscount       int
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
}

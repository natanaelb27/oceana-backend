package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMerchantPromotion struct {
	ID                         uint
	UserID                     uuid.UUID `gorm:"uniqueIndex:idx_user_merchant_promotion;"`
	User                       User
	ProductMerchantPromotionID uint `gorm:"uniqueIndex:idx_user_merchant_promotion;"`
	ProductMerchantPromotion   ProductMerchantPromotion
	Usage                      uint
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
	DeletedAt                  gorm.DeletedAt
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMerchantPromotion struct {
	ID                         uint
	UserID                     uuid.UUID
	User                       User
	ProductMerchantPromotionID uint
	ProductMerchantPromotion   ProductMerchantPromotion
	Usage                      uint
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
	DeletedAt                  gorm.DeletedAt
}

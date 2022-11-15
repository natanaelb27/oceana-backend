package models

import (
	"time"

	"gorm.io/gorm"
)

type MerchantPromotion struct {
	ID                    uint
	MerchantId            uint
	Merchant              Merchant
	Name                  string
	StartDate             time.Time
	EndDate               time.Time
	DiscountType          string
	Discount              uint
	MaximumDiscountAmount uint
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt
}

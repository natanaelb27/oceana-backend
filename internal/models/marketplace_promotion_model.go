package models

import (
	"time"

	"gorm.io/gorm"
)

type MarketplacePromotion struct {
	ID                     uint
	Name                   string
	StartDate              time.Time
	EndDate                time.Time
	Discount               uint
	MaxmimumDiscountAmount uint
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt
}

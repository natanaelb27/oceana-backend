package models

import (
	"time"

	"gorm.io/gorm"
)

type MarketplaceVoucher struct {
	ID                     uint
	Name                   string
	StartDate              time.Time
	EndDate                time.Time
	DiscountType           string
	Discount               int
	MinimumSpendAmount     int
	MaxmimumDiscountAmount int
	MaximumQuota           int
	CurrentQuota           int
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt
}

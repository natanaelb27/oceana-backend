package models

import (
	"time"

	"gorm.io/gorm"
)

type MarketplaceVoucher struct {
	ID                    uint
	BannerUrl             string
	Name                  string
	StartDatetime         time.Time
	EndDatetime           time.Time
	DiscountType          string
	Discount              uint
	MinimumSpendAmount    uint
	MaximumDiscountAmount uint
	MaximumQuota          uint
	CurrentQuota          uint
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt
}

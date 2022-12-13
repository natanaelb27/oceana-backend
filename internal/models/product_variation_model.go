package models

import (
	"time"
)

type ProductVariation struct {
	ID               uint
	ProductItemID    uint
	ProductItem      ProductItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductVariantID uint
	ProductVariant   ProductVariant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

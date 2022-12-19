package models

import (
	"time"
)

type ProductVariantMedia struct {
	ID               uint
	ProductVariantID uint
	ProductVariant   ProductVariant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MediaURL         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

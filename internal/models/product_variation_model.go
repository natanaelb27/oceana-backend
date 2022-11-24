package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariation struct {
	ID               uint
	ProductItemID    uint
	ProductItem      ProductItem
	ProductVariantID uint
	ProductVariant   ProductVariant
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

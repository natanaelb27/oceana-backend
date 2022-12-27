package models

import (
	"time"
)

type ProductVariation struct {
	ID               uint
	ProductItemID    uint           `gorm:"uniqueIndex:idx_product_item_id_product_variant_id;"`
	ProductItem      ProductItem    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductVariantID uint           `gorm:"uniqueIndex:idx_product_item_id_product_variant_id;"`
	ProductVariant   ProductVariant `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

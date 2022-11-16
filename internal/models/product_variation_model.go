package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariation struct {
	ID                      uint
	ProductItemID           uint
	ProductItem             ProductItem
	ProductVariationGroupID uint
	ProductVariationGroup   ProductVariationGroup
	Name                    string
	CreatedAt               time.Time
	UpdatedAt               time.Time
	DeletedAt               gorm.DeletedAt
}

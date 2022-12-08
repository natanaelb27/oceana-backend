package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariationSnapshot struct {
	ID                uint
	ProductSnapshotID uint `json:"product_snapshot_id"`
	ProductSnapshot   ProductSnapshot
	VariantType       string `json:"variant_type"`
	Variant           string `json:"variant"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

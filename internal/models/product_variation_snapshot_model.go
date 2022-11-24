package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariationSnapshot struct {
	ID                uint
	ProductSnapshotID uint
	ProductSnapshot   ProductSnapshot
	VariantType       string
	Variant           string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

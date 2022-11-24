package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariantMedia struct {
	ID               uint
	ProductVariantID uint
	ProductVariant   ProductVariant
	MediaURL         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

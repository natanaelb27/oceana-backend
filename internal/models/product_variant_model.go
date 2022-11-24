package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	ID                   uint
	ProductVariantTypeID uint
	ProductVariantType   ProductVariantType
	Name                 string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt
}

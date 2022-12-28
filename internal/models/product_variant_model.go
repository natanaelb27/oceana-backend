package models

import (
	"time"
)

type ProductVariant struct {
	ID                   uint
	ProductVariantTypeID uint
	ProductVariantType   ProductVariantType `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name                 string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

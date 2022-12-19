package models

import (
	"time"
)

type ProductVariantType struct {
	ID        uint
	ProductID uint
	Product   Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductVariationGroup struct {
	ID        uint
	ProductID uint
	Product   Product
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

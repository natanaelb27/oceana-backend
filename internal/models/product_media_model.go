package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductMedia struct {
	ID        uint
	ProductID uint
	Product   Product
	MediaURL  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

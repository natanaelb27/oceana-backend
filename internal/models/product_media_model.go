package models

import (
	"time"
)

type ProductMedia struct {
	ID        uint
	ProductID uint
	Product   Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MediaURL  string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductItemMedia struct {
	ID            uint
	ProductItemID uint
	ProductItem   ProductItem
	MediaURL      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

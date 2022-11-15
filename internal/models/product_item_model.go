package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductItem struct {
	ID            uint
	ProductID     uint
	Product       Product
	SKU           string
	Price         int
	StockQuantity int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

package models

import (
	"time"
)

type ProductItem struct {
	ID            uint
	ProductID     uint
	Product       Product
	SKU           string
	Price         uint
	StockQuantity uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

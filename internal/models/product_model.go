package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint
	MerchantID    uint
	Merchant      Merchant
	CategoryID    uint
	Category      Category
	Name          string
	Description   string
	IsHazardous   bool
	Weight        uint
	PackageSize   string
	Condition     string
	ProductStatus string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

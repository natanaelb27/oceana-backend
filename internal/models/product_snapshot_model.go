package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductSnapshot struct {
	ID          uint
	CategoryID  uint
	Category    Category
	ProductID   uint
	Product     Product
	Name        string
	Description string
	IsHazardous bool
	Weight      uint
	PackageSize string
	Condition   string
	SKU         string
	Price       uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductSnapshot struct {
	ID          uint
	CategoryID  uint
	Category    Category
	ProductID   int
	Product     Product
	Name        string
	Description string
	SKU         string
	Price       uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

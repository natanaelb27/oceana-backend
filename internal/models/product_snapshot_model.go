package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductSnapshot struct {
	ID          uint
	CategoryID  uint
	Category    Category
	Name        string
	Price       uint
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

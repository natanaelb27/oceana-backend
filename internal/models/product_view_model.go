package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductView struct {
	ID           uint
	UserId       uint
	User         User
	ProductId    uint
	Product      Product
	LastViewedAt time.Time
	Count        uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

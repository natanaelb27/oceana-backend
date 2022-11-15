package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductFavorite struct {
	ID         uint
	UserId     uint
	User       User
	ProductId  uint
	Product    Product
	IsFavorite bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

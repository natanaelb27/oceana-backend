package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductFavorite struct {
	ID         uint
	UserId     uuid.UUID
	User       User
	ProductId  uint
	Product    Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsFavorite bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

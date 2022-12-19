package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductView struct {
	ID           uint
	UserId       uuid.UUID
	User         User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductId    uint
	Product      Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LastViewedAt time.Time
	Count        uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

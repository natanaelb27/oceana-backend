package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID            uint
	UserID        uuid.UUID
	User          User
	ProductItemID uint
	ProductItem   ProductItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Quantity      uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID            uint
	UserID        uuid.UUID
	User          User
	ProductItemID uint
	ProductItem   ProductItem
	Quantity      int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

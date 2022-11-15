package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Merchant struct {
	ID        uint
	UserID    uuid.UUID
	User      User
	Name      string
	AddressID uuid.UUID
	Address   Address
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

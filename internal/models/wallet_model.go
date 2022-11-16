package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID                    uuid.UUID
	UserID                uuid.UUID
	User                  User
	Balance               uint
	Pin                   string
	IsActive              bool
	InvalidPinCount       int
	BlockedExpirationTime *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
	DeletedAt             gorm.DeletedAt
}

func (m *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

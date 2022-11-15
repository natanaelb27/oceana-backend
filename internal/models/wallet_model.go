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
	InvalidPinCount       uint
	BlockedExpirationTime time.Time
}

func (m *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

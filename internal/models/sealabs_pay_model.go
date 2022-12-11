package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SealabsPay struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	User       User
	CardNumber string
	ActiveDate time.Time
	NameOnCard string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

func (m *SealabsPay) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

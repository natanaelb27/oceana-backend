package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SealabsPay struct {
	ID            uuid.UUID
	UserId        uuid.UUID
	User          User
	AccountNumber string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (m *SealabsPay) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

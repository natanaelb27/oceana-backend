package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Email struct {
	ID      uuid.UUID
	Address string
}

func (m *Email) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

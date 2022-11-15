package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID
	FullName    string
	Username    string
	EmailID     uuid.UUID
	Email       Email
	Password    string
	PhoneNumber string
	Gender      string
	Birthdate   time.Time
	Photo       *string
	GoogleID    *string
	Role        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func (m *User) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New()
	return
}

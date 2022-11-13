package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

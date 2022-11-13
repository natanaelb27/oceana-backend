package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

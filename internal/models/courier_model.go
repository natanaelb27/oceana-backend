package models

import (
	"time"

	"gorm.io/gorm"
)

type Courier struct {
	ID        uint
	Code      string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

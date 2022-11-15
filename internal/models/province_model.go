package models

import (
	"time"

	"gorm.io/gorm"
)

type Province struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

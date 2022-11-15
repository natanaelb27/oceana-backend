package models

import (
	"time"

	"gorm.io/gorm"
)

type City struct {
	ID         uint
	Name       string
	ProvinceId uint
	Province   Province
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

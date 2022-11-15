package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductMediaSnapshot struct {
	ID                uint
	ProductSnapshotID uint
	ProductSnapshot   ProductSnapshot
	MediaURL          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductMediaSnapshot struct {
	ID                uint            `json:"id"`
	ProductSnapshotID uint            `json:"product_snapshot_id"`
	ProductSnapshot   ProductSnapshot `json:"product_snapshot"`
	MediaURL          string          `json:"media_url"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `json:"deleted_at"`
}

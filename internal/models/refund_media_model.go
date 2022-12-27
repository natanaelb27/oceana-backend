package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefundMedia struct {
	ID        uint
	RefundID  uuid.UUID
	MediaURL  string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

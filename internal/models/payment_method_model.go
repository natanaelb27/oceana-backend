package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentMethod struct {
	ID        uint
	Name      string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

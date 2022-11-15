package models

import (
	"time"

	"gorm.io/gorm"
)

type CourierService struct {
	ID          uint
	CourierID   uint
	Courier     Courier
	Code        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

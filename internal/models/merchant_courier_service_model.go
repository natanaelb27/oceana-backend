package models

import (
	"time"

	"gorm.io/gorm"
)

type MerchantCourierService struct {
	ID               uint
	CourierServiceId uint
	CourierService   CourierService
	MerchantID       uint
	Merchant         Merchant
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

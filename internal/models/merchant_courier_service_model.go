package models

import (
	"time"

	"gorm.io/gorm"
)

type MerchantCourierService struct {
	ID               uint
	CourierServiceId uint `gorm:"uniqueIndex:idx_merchant_courier_service;"`
	CourierService   CourierService
	MerchantID       uint `gorm:"uniqueIndex:idx_merchant_courier_service;"`
	Merchant         Merchant
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

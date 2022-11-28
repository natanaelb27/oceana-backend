package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductCourierService struct {
	ID                       uint
	MerchantCourierServiceID uint
	MerchantCourierService   MerchantCourierService
	ProductId                uint
	Product                  Product
	CreatedAt                time.Time
	UpdatedAt                time.Time
	DeletedAt                gorm.DeletedAt
}

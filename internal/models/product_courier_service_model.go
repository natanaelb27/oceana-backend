package models

import (
	"time"
)

type ProductCourierService struct {
	ID                       uint
	MerchantCourierServiceID uint
	MerchantCourierService   MerchantCourierService `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductId                uint
	Product                  Product
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

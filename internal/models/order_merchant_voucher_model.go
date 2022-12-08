package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderMerchantVoucher struct {
	ID                uint
	OrderID           uuid.UUID
	Order             Order
	MerchantVoucherID uint
	MerchantVoucher   MerchantVoucher
	TotalDiscount     float64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

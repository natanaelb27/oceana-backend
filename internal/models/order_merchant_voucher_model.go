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
	TotalDiscount     uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt
}

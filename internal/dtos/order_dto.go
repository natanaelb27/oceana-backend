package dtos

import (
	"time"

	"github.com/google/uuid"
)

type OrderList struct {
	TransactionID       uuid.UUID          `json:"transaction_id"`
	OrderID             uuid.UUID          `json:"order_id"`
	MerchantID          uint               `json:"merchant_id"`
	ShopName            string             `json:"shop_name"`
	OrderDetails        []*OrderDetailList `json:"order_details" gorm:"-"`
	TransactionDatetime time.Time          `json:"transaction_datetime"`
	OrderStatus         string             `json:"order_status"`
	TotalPrice          float64            `json:"total_price"`
}

// type OrderListDetails struct {
// 	TransactionID                  uuid.UUID          `json:"transaction_id"`
// 	OrderID                        uuid.UUID          `json:"oder_id"`
// 	OrderDetails                   []*OrderDetailList `json:"order_details"`
// 	PaymentMethod                  string             `json:"payment_method"`
// 	CourierDescription             string             `json:"courier_description"`
// 	TotalMerchantPromotionDiscount uint               `json:"total_merchant_promotion_discount"`
// 	MerchantVoucherDiscount        uint               `json:"merchant_voucher_discount"`
// 	MarketplacePromotionDiscount   uint               `json:"marketplace_promotion_discount"`
// 	MarketplaceVoucherDiscount     uint               `json:"marketplace_voucher_discount"`
// 	TotalPriceBeforeDiscount       uint               `json:"total_price_before_discount"`
// 	TotalPriceAfterDiscount        uint               `json:"total_price_after_discount"`
// }

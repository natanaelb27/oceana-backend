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

type OrderListDetails struct {
	TransactionID                     uuid.UUID               `json:"transaction_id"`
	UserID                            uuid.UUID               `json:"user_id"`
	OrderID                           uuid.UUID               `json:"oder_id"`
	MerchantID                        uint                    `json:"merchant_id"`
	ShopName                          string                  `json:"shop_name"`
	OrderDetails                      []*OrderDetailList      `json:"order_details" gorm:"-"`
	TransactionDatetime               time.Time               `json:"transaction_datetime"`
	OrderStatus                       string                  `json:"order_status"`
	CourierDescription                string                  `json:"courier_description"`
	Address                           *AddressSnapshotDetails `json:"address" gorm:"-"`
	PaymentMethod                     string                  `json:"payment_method"`
	TotalPriceBeforeDiscount          uint                    `json:"total_price_before_discount"`
	TotalMerchantPromotionDiscount    float64                 `json:"total_merchant_promotion_discount"`
	TotalMerchantVoucherDiscount      float64                 `json:"total_merchant_voucher_discount"`
	TotalMarketplacePromotionDiscount float64                 `json:"total_marketplace_promotion_discount"`
	TotalMarketplaceVoucherDiscount   float64                 `json:"total_marketplace_voucher_discount"`
	ShippingFee                       uint                    `json:"shipping_fee"`
	TotalPayment                      float64                 `json:"total_payment"`
}

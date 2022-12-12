package dtos

import (
	"time"

	"github.com/google/uuid"
)

type OrderList struct {
	TransactionID       uuid.UUID          `json:"transaction_id"`
	OrderID             uuid.UUID          `json:"order_id"`
	InvoiceNumber       string             `json:"invoice_number"`
	MerchantID          uint               `json:"merchant_id"`
	ShopName            string             `json:"shop_name"`
	UserID              uuid.UUID          `json:"user_id"`
	UserFullName        string             `json:"user_full_name"`
	OrderDetails        []*OrderDetailList `json:"order_details" gorm:"-"`
	TransactionDatetime time.Time          `json:"transaction_datetime"`
	OrderStatus         string             `json:"order_status"`
	TotalPrice          float64            `json:"total_price"`
}

type OrderListDetails struct {
	TransactionID                   uuid.UUID               `json:"transaction_id"`
	UserID                          uuid.UUID               `json:"user_id"`
	UserFullName                    string                  `json:"user_full_name"`
	OrderID                         uuid.UUID               `json:"oder_id"`
	InvoiceNumber                   string                  `json:"invoice_number"`
	MerchantID                      uint                    `json:"merchant_id"`
	ShopName                        string                  `json:"shop_name"`
	OrderDetails                    []*OrderDetailList      `json:"order_details" gorm:"-"`
	TransactionDatetime             time.Time               `json:"transaction_datetime"`
	OrderStatus                     string                  `json:"order_status"`
	CourierDescription              string                  `json:"courier_description"`
	OrderTrackingHistory            []*OrderTrackingHistory `json:"order_tracking_history" gorm:"-"`
	Address                         *AddressSnapshotDetails `json:"address" gorm:"-"`
	PaymentMethod                   string                  `json:"payment_method"`
	TotalPriceBeforeDiscount        uint                    `json:"total_price_before_discount"`
	TotalMerchantPromotionDiscount  float64                 `json:"total_merchant_promotion_discount"`
	TotalMerchantVoucherDiscount    float64                 `json:"total_merchant_voucher_discount"`
	TotalMarketplaceVoucherDiscount float64                 `json:"total_marketplace_voucher_discount"`
	ShippingFee                     uint                    `json:"shipping_fee"`
	TotalPayment                    float64                 `json:"total_payment"`
}

type OrderProcessOrder struct {
	CourierServiceID uint
	OrderID          uuid.UUID
}

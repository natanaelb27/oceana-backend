package dtos

import "github.com/google/uuid"

type ProductItem struct {
	ProductItemID uint    `json:"product_item_id" binding:"required"`
	Quantity      uint    `json:"quantity" binding:"required"`
	Message       *string `json:"message"`
}

type Product struct {
	ProductID                  uint           `json:"product_id" binding:"required"`
	ProductMerchantPromotionID *uint          `json:"product_merchant_promotion_id"`
	ProductItems               []*ProductItem `json:"product_items" binding:"required"`
}

type CartOrderItem struct {
	MerchantID        uint       `json:"merchant_id" binding:"required"`
	MerchantVoucherID *uint      `json:"merchant_voucher_id"`
	CourierServiceID  uint       `json:"courier_service_id" binding:"required"`
	ShippingFee       uint       `json:"shipping_fee" binding:"required"`
	Products          []*Product `json:"products" binding:"required"`
}

type PayCartWithWalletRequest struct {
	UserID               uuid.UUID        `json:"user_id"`
	CartOrderItems       []*CartOrderItem `json:"cart_order_items" binding:"required"`
	MarketplaceVoucherID *uint            `json:"marketplace_voucher_id"`
	AddressID            uuid.UUID        `json:"address_id" binding:"required"`
	Pin                  string           `json:"pin" binding:"required,len=6"`
	TotalPayment         float64          `json:"total_payment" binding:"required"`
}

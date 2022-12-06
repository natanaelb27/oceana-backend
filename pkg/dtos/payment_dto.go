package dtos

import "github.com/google/uuid"

type OrderCartItem struct {
	ProductItemID uint    `json:"product_item_id" binding:"required"`
	Quantity      *uint   `json:"quantity" binding:"required"`
	Message       *string `json:"message"`
}

type CartOrderItem struct {
	OrderCartItems      []OrderCartItem `json:"order_cart_items" binding:"required"`
	MerchantID          uint            `json:"merchant_id" binding:"required"`
	MerchantVoucherID   *uint           `json:"merchant_voucher_id"`
	MerchantPromotionID *uint           `json:"merchant_promotion_id"`
	CourierServiceID    uint            `json:"courier_service_id" binding:"required"`
}

type PayCartWithWalletRequest struct {
	UserID               uuid.UUID
	CartOrderItems       []CartOrderItem `json:"cart_order_items" binding:"required"`
	MarketplaceVoucherID *uint           `json:"marketplace_voucher_id"`
	AddressID            uint            `json:"address_id" binding:"required"`
	Pin                  string          `json:"pin" binding:"required,len=6"`
}

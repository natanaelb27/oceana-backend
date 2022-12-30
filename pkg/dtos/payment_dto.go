package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

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
	CartIDs              []*uint          `json:"cart_ids" binding:"required"`
	MarketplaceVoucherID *uint            `json:"marketplace_voucher_id"`
	AddressID            uuid.UUID        `json:"address_id" binding:"required"`
	Pin                  string           `json:"pin" binding:"required,len=6"`
	TotalPayment         float64          `json:"total_payment" binding:"required"`
}

type PaymentResponse struct {
	ID                        uuid.UUID `json:"id"`
	UserID                    uuid.UUID `json:"user_id"`
	TransactionDatetime       time.Time `json:"transaction_datetime"`
	TotalPrice                float64   `json:"total_price"`
	TotalPriceAfterDiscount   float64   `json:"total_price_after_discount"`
	TotalPriceWithShippingFee float64   `json:"total_price_with_shipping"`
	PaymentMethodID           uint      `json:"payment_method_id"`
	IsPaid                    bool      `json:"is_paid"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

func NewPaymentResponse(source *models.Transaction) (target *PaymentResponse) {
	target = &PaymentResponse{
		ID:                        source.ID,
		UserID:                    source.UserID,
		TransactionDatetime:       source.TransactionDatetime,
		TotalPrice:                source.TotalPrice,
		TotalPriceAfterDiscount:   source.TotalPriceAfterDiscount,
		TotalPriceWithShippingFee: source.TotalPriceWithShippingFee,
		PaymentMethodID:           source.PaymentMethodID,
		IsPaid:                    source.IsPaid,
		CreatedAt:                 source.CreatedAt,
		UpdatedAt:                 source.UpdatedAt,
	}
	return
}

type PayCartWithSealabsPayRequest struct {
	UserID               uuid.UUID        `json:"user_id"`
	CartOrderItems       []*CartOrderItem `json:"cart_order_items" binding:"required"`
	CartIDs              []*uint          `json:"cart_ids" binding:"required"`
	MarketplaceVoucherID *uint            `json:"marketplace_voucher_id"`
	AddressID            uuid.UUID        `json:"address_id" binding:"required"`
	CardNumber           string           `json:"card_number" binding:"required"`
	TotalPayment         float64          `json:"total_payment" binding:"required"`
}

type PayCartWithSealabsPaySuccessRequest struct {
	TransactionID   uuid.UUID `json:"transaction_id" binding:"required"`
	SealabsPayTxnID string    `json:"txn_id" binding:"required"`
	CardNumber      string    `json:"card_number" binding:"required"`
	Signature       string    `json:"signature" binding:"required"`
}

package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type RefundItems struct {
	RefundDetailID     uuid.UUID `json:"refund_detail_id"`
	OrderDetailID      uuid.UUID `json:"order_detail_id"`
	ProductImage       string    `json:"product_image"`
	ProductName        string    `json:"product_name"`
	ProductVariations  string    `json:"product_variations"`
	RefundQuantity     uint      `json:"refund_quantity"`
	UserRefundAmount   float64   `json:"user_refund_amount"`
	MerchantLostAmount float64   `json:"merchant_lost_amount"`
}

type RefundDetails struct {
	RefundID                                   uuid.UUID      `json:"refund_id"`
	OrderID                                    uuid.UUID      `json:"order_id"`
	InvoiceNumber                              string         `json:"invoice_number"`
	MerchantID                                 uint           `json:"merchant_id"`
	ShopName                                   string         `json:"shop_name"`
	UserID                                     uuid.UUID      `json:"user_id"`
	UserFullName                               string         `json:"user_full_name"`
	RefundRequestDatetime                      time.Time      `json:"refund_request_datetime"`
	RejectedDatetime                           *time.Time     `json:"rejected_datetime"`
	RefundStatus                               string         `json:"refund_status"`
	RefundReason                               string         `json:"refund_reason"`
	Images                                     pq.StringArray `json:"images" gorm:"type:text[]"`
	Videos                                     pq.StringArray `json:"videos" gorm:"type:text[]"`
	RefundItems                                []*RefundItems `json:"refund_items" gorm:"-"`
	TotalUserRefundAmount                      float64        `json:"total_user_refund_amount"`
	TotalMerchantLostAmount                    float64        `json:"total_merchant_lost_amount"`
	RefundSellerConfirmationExpirationDatetime *time.Time     `json:"refund_seller_confirmation_expiration_datetime"`
}

// type RefundDetailDetails struct {
// 	RefundID              uuid.UUID      `json:"refund_id"`
// 	OrderID               uuid.UUID      `json:"order_id"`
// 	InvoiceNumber         string         `json:"invoice_number"`
// 	MerchantID            uint           `json:"merchant_id"`
// 	ShopName              string         `json:"shop_name"`
// 	UserID                uuid.UUID      `json:"user_id"`
// 	UserFullName          string         `json:"user_full_name"`
// 	RefundRequestDatetime time.Time      `json:"refund_request_datetime"`
// 	RejectedDatetime      *time.Time     `json:"rejected_datetime"`
// 	RefundStatus          string         `json:"refund_status"`
// 	Images                pq.StringArray `json:"images"`
// 	Videos                pq.StringArray `json:"videos"`
// 	TotalUserRefundAmount float64        `json:"total_user_refund_amount"`
// 	TotalMerchantLost     float64        `json:"total_merchant_lost"`
// }

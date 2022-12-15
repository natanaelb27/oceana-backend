package dtos

import (
	"time"

	"github.com/google/uuid"
)

type RefundDetails struct {
	RefundID              uuid.UUID  `json:"refund_id"`
	OrderID               uuid.UUID  `json:"order_id"`
	InvoiceNumber         string     `json:"invoice_number"`
	MerchantID            uint       `json:"merchant_id"`
	ShopName              string     `json:"shop_name"`
	UserID                uuid.UUID  `json:"user_id"`
	UserFullName          string     `json:"user_full_name"`
	RefundRequestDatetime time.Time  `json:"refund_request_datetime"`
	RejectedDatetime      *time.Time `json:"rejected_datetime"`
	RefundStatus          string     `json:"refund_status"`
	TotalUserRefundAmount float64    `json:"total_user_refund_amount"`
	TotalMerchantLost     float64    `json:"total_merchant_lost"`
}

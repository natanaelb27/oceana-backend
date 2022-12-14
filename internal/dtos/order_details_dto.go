package dtos

import "github.com/google/uuid"

type OrderDetailList struct {
	OrderDetailID     uuid.UUID `json:"order_detail_id"`
	Message           string    `json:"message"`
	ProductImage      string    `json:"product_image"`
	ProductName       string    `json:"product_name"`
	ProductVariations string    `json:"product_variations"`
	Quantity          uint      `json:"quantity"`
	Price             uint      `json:"price"`
	TotalPrice        uint      `json:"total_price"`
}

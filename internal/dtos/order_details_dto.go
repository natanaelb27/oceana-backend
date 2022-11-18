package dtos

type OrderDetailList struct {
	OrderDetailID     string `json:"order_detail_id"`
	ProductImage      string `json:"product_image"`
	ProductName       string `json:"product_name"`
	ProductVariations string `json:"product_variations"`
	Quantity          uint   `json:"quantity"`
	Price             uint   `json:"price"`
	TotalPrice        uint   `json:"total_price"`
}

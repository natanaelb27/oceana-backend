package dtos

import (
	"time"
)

type ProductMerchantPromotion struct {
	ProductID             uint    `json:"product_id"`
	ProductName           string  `json:"product_name"`
	ProductImage          string  `json:"product_image"`
	PriceMax              uint    `json:"price_max"`
	PriceMaxAfterDiscount float64 `json:"price_max_after_discount"`
	PriceMin              uint    `json:"price_min"`
	PriceMinAfterDiscount float64 `json:"price_min_after_discount"`
	TotalStock            uint    `json:"total_stock"`
	DiscountType          string  `json:"discount_type"`
	Discount              uint    `json:"discount"`
	MaximumDiscountAmount uint    `json:"maximum_discount_amount"`
	MaximumQuota          int     `json:"maximum_quota"`
	MaximumUserLimit      int     `json:"maximum_user_limit"`
	CurrentQuota          uint    `json:"current_quota"`
}

type MerchantPromotionDetails struct {
	ID                       uint                        `json:"id"`
	MerchantID               uint                        `json:"merchant_id"`
	Name                     string                      `json:"name"`
	StartDatetime            time.Time                   `json:"start_datetime"`
	EndDatetime              time.Time                   `json:"end_datetime"`
	PromotionStatus          string                      `json:"promotion_status"`
	ProductMerchantPromotion []*ProductMerchantPromotion `json:"product_merchant_promotion"`
}

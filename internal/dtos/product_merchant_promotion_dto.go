package dtos

import "time"

type ProductMerchantPromotionDetails struct {
	MerchantPromotionID        uint      `json:"merchant_promotion_id"`
	MerchantID                 uint      `json:"merchant_id"`
	Name                       string    `json:"name"`
	StartDatetime              time.Time `json:"start_datetime"`
	EndDatetime                time.Time `json:"end_datetime"`
	ProductMerchantPromotionID uint      `json:"product_merchant_promotion_id"`
	ProductID                  uint      `json:"product_id"`
	DiscountType               string    `json:"discount_type"`
	Discount                   uint      `json:"discount"`
	MaximumDiscountAmount      uint      `json:"maximum_discount_amount"`
	MaximumQuota               int       `json:"maximum_quota"`
	MaximumUserLimit           int       `json:"maximum_user_limit"`
	CurrentQuota               uint      `json:"current_quota"`
	PromotionStatus            string    `json:"promotion_status"`
}

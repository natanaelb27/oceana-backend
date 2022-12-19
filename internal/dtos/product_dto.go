package dtos

import "github.com/lib/pq"

type ProductDetails struct {
	Name                  string                `json:"name"`
	MerchantID            uint                  `json:"merchant_id"`
	IsHazardous           bool                  `json:"is_hazardous"`
	Weight                uint                  `json:"weight"`
	PackageSize           string                `json:"package_size"`
	Condition             string                `json:"condition"`
	CategoryID            uint                  `json:"category_id"`
	Category              string                `json:"category"`
	Images                pq.StringArray        `json:"images" gorm:"type:text[]"`
	Videos                pq.StringArray        `json:"videos" gorm:"type:text[]"`
	Description           string                `json:"description"`
	TotalStock            uint                  `json:"total_stock"`
	LikeCount             uint                  `json:"like_count"`
	TotalPurchased        uint                  `json:"total_purchased"`
	ProductReviewDetails  *ProductReviewDetails `json:"product_review_details"`
	DiscountType          *string               `json:"discount_type"`
	Discount              *uint                 `json:"discount"`
	Price                 float64               `json:"price"`
	PriceAfterDiscount    *float64              `json:"price_before_discount"`
	PriceMax              float64               `json:"price_max"`
	PriceMaxAfterDiscount *float64              `json:"price_max_before_discount"`
	PriceMin              float64               `json:"price_min"`
	PriceMinAfterDiscount *float64              `json:"price_min_before_discount"`
	ProductModel          []*ProductModel       `json:"product_model" gorm:"-"`
	ProductVariation      []*ProductVariation   `json:"product_variation" gorm:"-"`
}

type SelectedProductPurchaseQuantities struct {
	Price                      float64 `json:"price"`
	Stock                      uint    `json:"stock"`
	MaxPurchaseQuantity        uint    `json:"max_purchase_quantity"`
	UserMaxPurchaseQuantity    uint    `json:"user_max_purchase_quantity"`
	ProductMerchantPromotionID *uint   `json:"product_merchant_promotion_id"`
}

type AvailableProducts struct {
	ProductID    uint   `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductImage string `json:"product_image"`
	PriceMax     uint   `json:"price_max"`
	PriceMin     uint   `json:"price_min"`
	TotalStock   uint   `json:"total_stock"`
}

package dtos

import "github.com/lib/pq"

type ProductDetails struct {
	Name        string         `json:"name"`
	IsHazardous bool           `json:"is_hazardous"`
	Weight      uint           `json:"weight"`
	PackageSize string         `json:"package_size"`
	Condition   string         `json:"condition"`
	Category    string         `json:"category"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
	// ProductMedia          []*ProductMedia     `json:"product_media" gorm:"-"`
	Description           string              `json:"description"`
	TotalStock            uint                `json:"total_stock"`
	LikeCount             uint                `json:"like_count"`
	RatingCount           uint                `json:"rating_count"`
	RatingStar            float64             `json:"rating_star"`
	DiscountType          string              `json:"discount_type"`
	Discount              uint                `json:"discount"`
	Price                 float64             `json:"price"`
	PriceAfterDiscount    float64             `json:"price_before_discount"`
	PriceMax              float64             `json:"price_max"`
	PriceMaxAfterDiscount float64             `json:"price_max_before_discount"`
	PriceMin              float64             `json:"price_min"`
	PriceMinAfterDiscount float64             `json:"price_min_before_discount"`
	ProductModel          []*ProductModel     `json:"product_model" gorm:"-"`
	ProductVariation      []*ProductVariation `json:"product_variation" gorm:"-"`
}

package dtos

type ProductSpecification struct {
	IsHazardous bool
	Weight      uint
	PackageSize string
	Condition   string
	Category    string
}

type ProductDetails struct {
	Name                   string                `json:"name"`
	ProductSpecification   *ProductSpecification `json:"product_specification" gorm:"-"`
	ProductMedia           []*ProductMedia       `json:"product_media"`
	Description            string                `json:"description"`
	TotalStock             uint                  `json:"total_stock"`
	TotalLike              uint                  `json:"total_like"`
	TotalRating            uint                  `json:"total_rating"`
	RatingStar             float64               `json:"rating_star"`
	DiscountType           string                `json:"discount_type"`
	Discount               uint                  `json:"discount"`
	Price                  float64               `json:"price"`
	PriceBeforeDiscount    float64               `json:"price_before_discount"`
	PriceMax               float64               `json:"price_max"`
	PriceMaxBeforeDiscount float64               `json:"price_max_before_discount"`
	PriceMin               float64               `json:"price_min"`
	PriceMinBeforeDiscount float64               `json:"price_min_before_discount"`
	ProductModel           []*ProductModel       `json:"product_model"`
	ProductVariation       []*ProductVariation   `json:"product_variation"`
}

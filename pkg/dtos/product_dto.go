package dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"github.com/lib/pq"
)

type ProductDetailsGetByIDResponse struct {
	Name        string `json:"name"`
	IsHazardous bool   `json:"is_hazardous"`
	Weight      uint   `json:"weight"`
	PackageSize string `json:"package_size"`
	Condition   string `json:"condition"`
	Category    string `json:"category"`
	// ProductMedia          []*dtos.ProductMedia     `json:"product_media"`
	Images                pq.StringArray           `json:"images" gorm:"type:text[]"`
	Description           string                   `json:"description"`
	TotalStock            uint                     `json:"total_stock"`
	LikeCount             uint                     `json:"like_count"`
	RatingCount           uint                     `json:"rating_count"`
	RatingStar            float64                  `json:"rating_star"`
	DiscountType          string                   `json:"discount_type"`
	Discount              uint                     `json:"discount"`
	Price                 float64                  `json:"price"`
	PriceAfterDiscount    float64                  `json:"price_after_discount"`
	PriceMax              float64                  `json:"price_max"`
	PriceMaxAfterDiscount float64                  `json:"price_max_after_discount"`
	PriceMin              float64                  `json:"price_min"`
	PriceMinAfterDiscount float64                  `json:"price_min_after_discount"`
	ProductModel          []*dtos.ProductModel     `json:"product_model"`
	ProductVariation      []*dtos.ProductVariation `json:"product_variation"`
}

func NewProductDetailsGetByIDResponse(source *dtos.ProductDetails) (target *ProductDetailsGetByIDResponse) {
	target = &ProductDetailsGetByIDResponse{
		Name:                  source.Name,
		IsHazardous:           source.IsHazardous,
		Weight:                source.Weight,
		PackageSize:           source.PackageSize,
		Condition:             source.Condition,
		Category:              source.Category,
		Images:                source.Images,
		Description:           source.Description,
		TotalStock:            source.TotalStock,
		LikeCount:             source.LikeCount,
		RatingCount:           source.RatingCount,
		RatingStar:            source.RatingStar,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		Price:                 source.Price,
		PriceAfterDiscount:    source.PriceAfterDiscount,
		PriceMax:              source.PriceMax,
		PriceMaxAfterDiscount: source.PriceMaxAfterDiscount,
		PriceMin:              source.PriceMin,
		PriceMinAfterDiscount: source.PriceMinAfterDiscount,
		ProductModel:          source.ProductModel,
		ProductVariation:      source.ProductVariation,
	}

	return
}

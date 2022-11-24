package dtos

import (
	"strconv"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type ProductDetailsGetByIDResponse struct {
	Name                  string                   `json:"name"`
	IsHazardous           bool                     `json:"is_hazardous"`
	Weight                uint                     `json:"weight"`
	PackageSize           string                   `json:"package_size"`
	Condition             string                   `json:"condition"`
	Category              string                   `json:"category"`
	Images                pq.StringArray           `json:"images" gorm:"type:text[]"`
	Description           string                   `json:"description"`
	TotalStock            uint                     `json:"total_stock"`
	LikeCount             uint                     `json:"like_count"`
	RatingCount           uint                     `json:"rating_count"`
	RatingStar            float64                  `json:"rating_star"`
	DiscountType          *string                  `json:"discount_type"`
	Discount              *uint                    `json:"discount"`
	Price                 float64                  `json:"price"`
	PriceAfterDiscount    *float64                 `json:"price_after_discount"`
	PriceMax              float64                  `json:"price_max"`
	PriceMaxAfterDiscount *float64                 `json:"price_max_after_discount"`
	PriceMin              float64                  `json:"price_min"`
	PriceMinAfterDiscount *float64                 `json:"price_min_after_discount"`
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

type PurchaseQuantitiesQuery struct {
	ProductID     uint `json:"product_id"`
	ProductItemID uint `json:"product_item_id"`
	Quantity      uint `json:"quantity"`
}

func NewPurchaseQuantitiesQuery(c *gin.Context) *PurchaseQuantitiesQuery {
	var query PurchaseQuantitiesQuery

	productID := c.Query("product_id")
	productItemID := c.Query("product_item_id")
	quantity := c.Query("quantity")

	productIDInt, err := strconv.Atoi(productID)
	if err != nil || productIDInt <= 0 {
		query.ProductID = 0
	}

	query.ProductID = uint(productIDInt)

	productItemIDInt, err := strconv.Atoi(productItemID)
	if err != nil || productItemIDInt <= 0 {
		query.ProductID = 0
	}

	query.ProductItemID = uint(productItemIDInt)

	quantityInt, err := strconv.Atoi(quantity)
	if err != nil || quantityInt <= 0 {
		query.ProductID = 0
	}

	query.Quantity = uint(quantityInt)

	return &query

}

type GetPurchaseQuantitiesByIDResponse struct {
	Price                      uint  `json:"price"`
	Stock                      uint  `json:"stock"`
	MaxPurchaseQuantity        uint  `json:"max_purchase_quantity"`
	UserMaxPurchaseQuantity    uint  `json:"user_max_purchase_quantity"`
	ProductMerchantPromotionID *uint `json:"product_merchant_promotion_id"`
}

func NewGetPurchaseQuantitiesByIDResponse(source *dtos.SelectedProductPurchaseQuantities) (target *GetPurchaseQuantitiesByIDResponse) {
	target = &GetPurchaseQuantitiesByIDResponse{
		Price:                      source.Price,
		Stock:                      source.Stock,
		MaxPurchaseQuantity:        source.MaxPurchaseQuantity,
		UserMaxPurchaseQuantity:    source.UserMaxPurchaseQuantity,
		ProductMerchantPromotionID: source.ProductMerchantPromotionID,
	}

	return
}

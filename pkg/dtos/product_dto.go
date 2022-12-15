package dtos

import (
	"strconv"
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/dtos"
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type ProductDetailsGetByIDResponse struct {
	Name                  string                     `json:"name"`
	MerchantID            uint                       `json:"merchant_id"`
	IsHazardous           bool                       `json:"is_hazardous"`
	Weight                uint                       `json:"weight"`
	PackageSize           string                     `json:"package_size"`
	Condition             string                     `json:"condition"`
	CategoryID            uint                       `json:"category_id"`
	Category              string                     `json:"category"`
	Images                pq.StringArray             `json:"images" gorm:"type:text[]"`
	Videos                pq.StringArray             `json:"videos" gorm:"type:text[]"`
	Description           string                     `json:"description"`
	TotalStock            uint                       `json:"total_stock"`
	LikeCount             uint                       `json:"like_count"`
	ProductReviewDetails  *dtos.ProductReviewDetails `json:"product_review_details"`
	DiscountType          *string                    `json:"discount_type"`
	Discount              *uint                      `json:"discount"`
	Price                 float64                    `json:"price"`
	PriceAfterDiscount    *float64                   `json:"price_after_discount"`
	PriceMax              float64                    `json:"price_max"`
	PriceMaxAfterDiscount *float64                   `json:"price_max_after_discount"`
	PriceMin              float64                    `json:"price_min"`
	PriceMinAfterDiscount *float64                   `json:"price_min_after_discount"`
	ProductModel          []*dtos.ProductModel       `json:"product_model"`
	ProductVariation      []*dtos.ProductVariation   `json:"product_variation"`
}

func NewProductDetailsGetByIDResponse(source *dtos.ProductDetails) (target *ProductDetailsGetByIDResponse) {
	target = &ProductDetailsGetByIDResponse{
		Name:                  source.Name,
		MerchantID:            source.MerchantID,
		IsHazardous:           source.IsHazardous,
		Weight:                source.Weight,
		PackageSize:           source.PackageSize,
		Condition:             source.Condition,
		CategoryID:            source.CategoryID,
		Category:              source.Category,
		Images:                source.Images,
		Videos:                source.Videos,
		Description:           source.Description,
		TotalStock:            source.TotalStock,
		LikeCount:             source.LikeCount,
		ProductReviewDetails:  source.ProductReviewDetails,
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
	ProductID     uint `json:"product-id"`
	ProductItemID uint `json:"product-item-id"`
	Quantity      uint `json:"quantity"`
}

func NewPurchaseQuantitiesQuery(c *gin.Context) *PurchaseQuantitiesQuery {
	var query PurchaseQuantitiesQuery

	productID := c.Query("product-id")
	productItemID := c.Query("product-item-id")
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
	if err != nil || quantityInt < 0 {
		query.ProductID = 0
	}

	query.Quantity = uint(quantityInt)

	return &query

}

type GetPurchaseQuantitiesByIDResponse struct {
	Price                      float64 `json:"price"`
	Stock                      uint    `json:"stock"`
	MaxPurchaseQuantity        uint    `json:"max_purchase_quantity"`
	UserMaxPurchaseQuantity    uint    `json:"user_max_purchase_quantity"`
	ProductMerchantPromotionID *uint   `json:"product_merchant_promotion_id"`
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

type ProductCreateRequest struct {
	CategoryID  uint   `json:"category_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsHazardous bool   `json:"is_hazardous"`
	Weight      uint   `json:"weight" binding:"required"`
	PackageSize string `json:"package_size" binding:"required"`
	Condition   string `json:"condition" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

func (source *ProductCreateRequest) ToModel(merchantID uint) (target *models.Product) {
	target = &models.Product{
		MerchantID:  merchantID,
		CategoryID:  source.CategoryID,
		Name:        source.Name,
		Description: source.Description,
		IsHazardous: source.IsHazardous,
		Weight:      source.Weight,
		PackageSize: source.PackageSize,
		Condition:   source.Description,
		Status:      source.Status,
	}

	return
}

type ProductCreateResponse struct {
	ID          uint      `json:"id"`
	MerchantID  uint      `json:"merchant_id"`
	CategoryID  uint      `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsHazardous bool      `json:"is_hazardous"`
	Weight      uint      `json:"weight"`
	PackageSize string    `json:"package_size"`
	Condition   string    `json:"condition"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProductCreateResponse(source *models.Product) (target *ProductCreateResponse) {
	target = &ProductCreateResponse{
		ID:          source.ID,
		MerchantID:  source.MerchantID,
		CategoryID:  source.CategoryID,
		Name:        source.Name,
		Description: source.Description,
		IsHazardous: source.IsHazardous,
		Weight:      source.Weight,
		PackageSize: source.PackageSize,
		Condition:   source.Condition,
		Status:      source.Status,
		CreatedAt:   source.CreatedAt,
		UpdatedAt:   source.UpdatedAt,
	}

	return
}

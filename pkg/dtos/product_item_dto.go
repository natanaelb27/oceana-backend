package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
)

type ProductItemCreateRequest struct {
	SKU           string `json:"sku" binding:"required"`
	Price         uint   `json:"price" binding:"required"`
	StockQuantity uint   `json:"stock_quantity" binding:"required"`
}

func (source *ProductItemCreateRequest) ToModel(productID uint) (target *models.ProductItem) {
	target = &models.ProductItem{
		ProductID:     productID,
		SKU:           source.SKU,
		Price:         source.Price,
		StockQuantity: source.StockQuantity,
	}

	return
}

type ProductItemCreateResponse struct {
	ID            uint      `json:"id"`
	ProductID     uint      `json:"product_id"`
	SKU           string    `json:"sku"`
	Price         uint      `json:"price"`
	StockQuantity uint      `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewProductItemCreateResponse(source *models.ProductItem) (target *ProductItemCreateResponse) {
	target = &ProductItemCreateResponse{
		ID:            source.ID,
		ProductID:     source.ProductID,
		SKU:           source.SKU,
		Price:         source.Price,
		StockQuantity: source.StockQuantity,
		CreatedAt:     source.CreatedAt,
		UpdatedAt:     source.UpdatedAt,
	}

	return
}

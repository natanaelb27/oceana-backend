package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
)

type ProductCreateRequest struct {
	CategoryID    uint   `json:"category_id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description" binding:"required"`
	IsHazardous   bool   `json:"is_hazardous"`
	Weight        uint   `json:"weight" binding:"required"`
	PackageSize   string `json:"package_size" binding:"required"`
	Condition     string `json:"condition" binding:"required"`
	ProductStatus string `json:"product_status" binding:"required"`
}

func (source *ProductCreateRequest) ToModel(merchantID uint) (target *models.Product) {
	target = &models.Product{
		MerchantID:    merchantID,
		CategoryID:    source.CategoryID,
		Name:          source.Name,
		Description:   source.Description,
		IsHazardous:   source.IsHazardous,
		Weight:        source.Weight,
		PackageSize:   source.PackageSize,
		Condition:     source.Description,
		ProductStatus: source.ProductStatus,
	}

	return
}

type ProductCreateResponse struct {
	ID            uint      `json:"id"`
	MerchantID    uint      `json:"merchant_id"`
	CategoryID    uint      `json:"category_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	IsHazardous   bool      `json:"is_hazardous"`
	Weight        uint      `json:"weight"`
	PackageSize   string    `json:"package_size"`
	Condition     string    `json:"condition"`
	ProductStatus string    `json:"product_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewProductCreateResponse(source *models.Product) (target *ProductCreateResponse) {
	target = &ProductCreateResponse{
		ID:            source.ID,
		MerchantID:    source.MerchantID,
		CategoryID:    source.CategoryID,
		Name:          source.Name,
		Description:   source.Description,
		IsHazardous:   source.IsHazardous,
		Weight:        source.Weight,
		PackageSize:   source.PackageSize,
		Condition:     source.Condition,
		ProductStatus: source.ProductStatus,
		CreatedAt:     source.CreatedAt,
		UpdatedAt:     source.UpdatedAt,
	}

	return
}

package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
)

type ProductMerchantPromotionRequest struct {
	ProductID             uint   `json:"product_id" binding:"required"`
	DiscountType          string `json:"discount_type" binding:"required"`
	Discount              uint   `json:"discount" binding:"required"`
	MaximumDiscountAmount uint   `json:"maximum_discount_amount"`
	MaximumQuota          int    `json:"maximum_quota" binding:"required"`
	MaximumUserLimit      int    `json:"maximum_user_limit" binding:"required"`
}

func (source *ProductMerchantPromotionRequest) ToModel(merchantPromotionId uint) (target *models.ProductMerchantPromotion) {
	target = &models.ProductMerchantPromotion{
		MerchantPromotionId:   merchantPromotionId,
		ProductID:             source.ProductID,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
		MaximumUserLimit:      source.MaximumUserLimit,
		CurrentQuota:          0,
	}

	return
}

type ProductMerchantPromotionResponse struct {
	ID                    uint      `json:"id"`
	MerchantPromotionId   uint      `json:"merchant_promotion_model_id"`
	ProductID             uint      `json:"product_id"`
	DiscountType          string    `json:"discount_type"`
	Discount              uint      `json:"discount"`
	MaximumDiscountAmount uint      `json:"maximum_discount_amount"`
	MaximumQuota          int       `json:"maximum_quota"`
	MaximumUserLimit      int       `json:"maximum_user_limit"`
	CurrentQuota          uint      `json:"current_quota"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func NewProductMerchantPromotionResponse(source *models.ProductMerchantPromotion) (target *ProductMerchantPromotionResponse) {
	target = &ProductMerchantPromotionResponse{
		ID:                    source.ID,
		MerchantPromotionId:   source.MerchantPromotionId,
		ProductID:             source.ProductID,
		DiscountType:          source.DiscountType,
		Discount:              source.Discount,
		MaximumDiscountAmount: source.MaximumDiscountAmount,
		MaximumQuota:          source.MaximumQuota,
		MaximumUserLimit:      source.MaximumUserLimit,
		CurrentQuota:          source.CurrentQuota,
		CreatedAt:             source.CreatedAt,
		UpdatedAt:             source.UpdatedAt,
	}

	return
}

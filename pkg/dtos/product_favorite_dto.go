package dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type GetProductFavoriteResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	ProductID  uint      `json:"product_id"`
	IsFavorite bool      `json:"is_favorite"`
}

func NewGetProductFavoriteResponse(source *models.ProductFavorite) (target *GetProductFavoriteResponse) {
	target = &GetProductFavoriteResponse{
		UserID:     source.UserId,
		ProductID:  source.ProductId,
		IsFavorite: source.IsFavorite,
	}

	return
}

type DeleteFavoriteProductsRequest struct {
	ProductIDs []*uint `json:"product_ids" binding:"required"`
}

type UpdateProductFavoriteResponse struct {
	UserID     uuid.UUID `json:"user_id"`
	ProductID  uint      `json:"product_id"`
	IsFavorite bool      `json:"is_favorite"`
}

func NewBatchUpdateProductFavoriteResponse(source []*models.ProductFavorite) (target []*UpdateProductFavoriteResponse) {
	for _, result := range source {
		target = append(target, &UpdateProductFavoriteResponse{
			UserID:     result.UserId,
			ProductID:  result.ProductId,
			IsFavorite: result.IsFavorite,
		})
	}

	return
}

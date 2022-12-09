package dtos

import (
	"time"

	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type GetProductViewResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	ProductID    uint      `json:"product_id"`
	LastViewedAt time.Time `json:"last_viewed_at"`
	Count        uint      `json:"count"`
}

func NewGetProductViewResponse(source *models.ProductView) (target *GetProductViewResponse) {
	target = &GetProductViewResponse{
		UserID:       source.UserId,
		ProductID:    source.ProductId,
		LastViewedAt: source.LastViewedAt,
		Count:        source.Count,
	}

	return
}

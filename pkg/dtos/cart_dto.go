package dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type CartUpsertRequest struct {
	ProductItemID uint  `json:"product_item_id" binding:"required,numeric"`
	Quantity      *uint `json:"quantity" binding:"required"`
}

type CartDeleteResponse struct {
	ID uint `json:"id"`
}

func NewCartDeleteResponse(source *models.Cart) (target *CartDeleteResponse) {
	target = &CartDeleteResponse{
		ID: source.ID,
	}
	return
}

type CartResponse struct {
	ID            uint      `json:"id"`
	ProductItemID uint      `json:"product_item_id"`
	UserID        uuid.UUID `json:"user_id"`
	Quantity      uint      `json:"quantity"`
}

func NewCartResponse(source *models.Cart) (target *CartResponse) {
	target = &CartResponse{
		ID:            source.ID,
		ProductItemID: source.ProductItemID,
		UserID:        source.UserID,
		Quantity:      source.Quantity,
	}
	return
}

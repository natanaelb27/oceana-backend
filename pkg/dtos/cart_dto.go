package dtos

import (
	"github.com/google/uuid"
	in_cart_dtos "github.com/natanaelb27/oceana-backend/internal/dtos/cart_dtos"
	"github.com/natanaelb27/oceana-backend/internal/models"
)

type CartUpsertRequest struct {
	ProductItemID uint    `json:"product_item_id" binding:"required,numeric"`
	Quantity      *uint   `json:"quantity" binding:"required"`
	Action        *string `json:"action"`
}

func (source *CartUpsertRequest) CartUpsertRequestToInternalCartRequest() (target *in_cart_dtos.CartRequest) {
	target = &in_cart_dtos.CartRequest{
		ProductItemID: source.ProductItemID,
		Quantity:      source.Quantity,
		Action:        source.Action,
	}
	return
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

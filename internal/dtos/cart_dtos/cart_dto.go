package cart_dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"git.garena.com/azmi-adhani/group-project-backend/pkg/dtos"
	"github.com/google/uuid"
)

type CartRequest struct {
	UserID        uuid.UUID
	ProductItemID uint
	Quantity      *uint
}

func NewInternalCartRequest(source *dtos.CartUpsertRequest) (target *CartRequest) {
	target = &CartRequest{
		ProductItemID: source.ProductItemID,
		Quantity:      source.Quantity,
	}
	return
}

func (source *CartRequest) ToModel() (target *models.Cart) {
	target = &models.Cart{
		UserID:        source.UserID,
		ProductItemID: source.ProductItemID,
		Quantity:      *source.Quantity,
	}
	return
}

type CartResult struct {
	UserID        uuid.UUID `json:"user_id"`
	ProductItemID uint      `json:"product_item_id"`
	Quantity      uint      `json:"name"`
}

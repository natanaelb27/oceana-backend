package cart_dtos

import (
	"git.garena.com/azmi-adhani/group-project-backend/internal/models"
	"github.com/google/uuid"
)

type CartRequest struct {
	UserID        uuid.UUID
	ProductItemID uint
	Quantity      *uint
}

func (source *CartRequest) ToModel() (target *models.Cart) {
	target = &models.Cart{
		UserID:        source.UserID,
		ProductItemID: source.ProductItemID,
		Quantity:      *source.Quantity,
	}
	return
}

type Product struct {
	CartsID            uint   `json:"carts_id"`
	ProductID          uint   `json:"product_id"`
	ProductItemID      uint   `json:"product_item_id"`
	ProductName        string `json:"product_name"`
	Quantity           uint   `json:"quantity"`
	Price              uint   `json:"price"`
	PriceAfterDiscount uint   `json:"price_after_discount"`
}

type Merchant struct {
	MerchantID uint       `json:"merchant_id"`
	Name       string     `json:"name"`
	Products   []*Product `json:"products" gorm:"-"`
}

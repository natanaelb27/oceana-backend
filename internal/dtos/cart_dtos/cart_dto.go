package cart_dtos

import (
	"github.com/google/uuid"
	"github.com/natanaelb27/oceana-backend/internal/models"
)

type CartRequest struct {
	UserID        uuid.UUID
	ProductItemID uint
	Quantity      *uint
	Action        *string
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
	CartsID                    uint     `json:"carts_id"`
	ProductID                  uint     `json:"product_id"`
	ProductItemID              uint     `json:"product_item_id"`
	ProductName                string   `json:"product_name"`
	ProductWeight              uint     `json:"product_weight"`
	ProductVariation           string   `json:"product_variation"`
	ProductMediaURL            string   `json:"product_media_url"`
	Quantity                   uint     `json:"quantity"`
	StockQuantity              uint     `json:"stock_quantity"`
	Price                      uint     `json:"price"`
	PriceAfterDiscount         *float64 `json:"price_after_discount"`
	ProductMerchantPromotionID uint     `json:"product_merchant_promotion_id"`
	PromoUserQuota             uint     `json:"promo_user_quota"`
	PromoGlobalQuota           uint     `json:"promo_global_quota"`
}

type Merchant struct {
	MerchantID       uint       `json:"merchant_id"`
	Name             string     `json:"name"`
	MerchantProvince string     `json:"merchant_province"`
	CityID           uint       `json:"city_id"`
	Products         []*Product `json:"products" gorm:"-"`
}

type OrderCartItem struct {
	ProductItemID      uint
	Quantity           *uint
	Message            *string
	Price              uint
	PriceAfterDiscount uint
}

type CartOrderItem struct {
	OrderCartItems          []OrderCartItem
	MerchantID              uint
	MerchantVoucherID       *uint
	MerchantPromotionID     *uint
	CourierServiceID        uint
	TotalPrice              uint
	TotalPriceAfterDiscount uint
}

type CartPayment struct {
	UserID                  uuid.UUID
	CartOrderItems          []CartOrderItem
	MarketplaceVoucherID    *uint
	AddressID               uint
	TotalPrice              uint
	TotalPriceAfterDiscount uint
}

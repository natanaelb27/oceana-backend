package dtos

import (
	"github.com/lib/pq"
)

type ProductModel struct {
	ProductItemID      uint          `json:"product_item_id"`
	SKU                string        `json:"sku"`
	StockQuantity      uint          `json:"stock_quantity"`
	Price              uint          `json:"price"`
	PriceAfterDiscount *float64      `json:"price_after_discount"`
	ProductVariations  string        `json:"product_variations"`
	VariationIndex     pq.Int64Array `json:"variation_index" gorm:"type:integer[]"`
}

type ProductItemDetails struct {
	ProductItemID uint   `json:"product_item_id"`
	SKU           string `json:"sku"`
	StockQuantity uint   `json:"stock_quantity"`
	Price         uint   `json:"price"`
	CategoryID    uint   `json:"category_id"`
	ProductID     uint   `json:"product_id"`
	MerchantID    uint   `json:"merchant_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	IsHazardous   bool   `json:"is_hazardous"`
	Weight        uint   `json:"weight"`
	PackageSize   string `json:"package_size"`
	Condition     string `json:"condition"`
}

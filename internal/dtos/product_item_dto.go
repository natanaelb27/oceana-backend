package dtos

import "github.com/lib/pq"

// type ProductModel struct {
// 	ProductItemID       uint           `json:"product_item_id"`
// 	SKU                 string         `json:"sku"`
// 	StockQuantity       uint           `json:"stock_quantity"`
// 	Price               uint           `json:"price"`
// 	PriceBeforeDiscount uint           `json:"price_before_discount"`
// 	Variation           pq.StringArray `json:"variation" gorm:"type:text[]"`
// }

type ProductModel struct {
	ProductItemID       uint          `json:"product_item_id"`
	SKU                 string        `json:"sku"`
	StockQuantity       uint          `json:"stock_quantity"`
	Price               uint          `json:"price"`
	PriceBeforeDiscount uint          `json:"price_before_discount"`
	VariationIndex      pq.Int64Array `json:"variation_index" gorm:"type:integer[]"`
}

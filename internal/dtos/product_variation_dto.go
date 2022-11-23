package dtos

import "github.com/lib/pq"

type Options struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type ProductVariation struct {
	Name string `json:"name"`
	// Options []*Options `json:"options"`
	// ProductVariationMedia []*ProductVariationMedia `json:"product_variation_media"`
	Options      pq.StringArray `json:"options" gorm:"type:text[]"`
	Image        pq.StringArray `json:"image" gorm:"type:text[]"`
	SummedStocks pq.Int64Array  `json:"summed_stocks" gorm:"type:integer[]"`
}

package dtos

import "github.com/lib/pq"

type ProductVariation struct {
	Name         string         `json:"name"`
	Options      pq.StringArray `json:"options" gorm:"type:text[]"`
	SummedStocks pq.Int64Array  `json:"summed_stocks" gorm:"type:integer[]"`
}

type ProductItemVariation struct {
	Variation pq.StringArray `json:"variation" gorm:"type:text[]"`
}

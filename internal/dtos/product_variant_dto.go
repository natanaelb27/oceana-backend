package dtos

import "github.com/lib/pq"

type ProductVariationTypeListOfVariants struct {
	Variants pq.Int64Array `json:"variants" gorm:"type:integer[]"`
}
